const MEMORY = {
    token: null,  // JWT
    dek: null,    // Decryption Key
    socket: null, 
    term: null,
    username: null
};

const API_URL = "http://localhost:8080";
const WS_URL = "ws://localhost:8080";

window.onload = () => {
    restoreSession();
};

function saveSession(data) {
    sessionStorage.setItem("pam_session", JSON.stringify({
        token: data.token,
        dek: data.dek,
        username: data.username
    }));
}

function restoreSession() {
    const saved = sessionStorage.getItem("pam_session");
    if (saved) {
        const data = JSON.parse(saved);
        
        MEMORY.token = data.token;
        MEMORY.dek = data.dek;
        MEMORY.username = data.username;

        showScreen("server-screen");
        document.getElementById("welcome-msg").innerText = `Welcome, ${data.username}`;
        
        if (data.username.toLowerCase() === "admin") {
            document.getElementById("admin-btn").classList.remove("hidden");
        }
        
        loadServers(data.username);
    }
}

function clearSession() {
    sessionStorage.removeItem("pam_session");
    location.reload();
}

function showScreen(id) {
    ['login-screen', 'server-screen', 'admin-screen', 'terminal-screen'].forEach(s => {
        document.getElementById(s).classList.add('hidden');
    });
    document.getElementById(id).classList.remove('hidden');
}

function setStatus(msg, isError = false) {
    const el = document.getElementById("admin-status");
    el.innerText = msg;
    el.style.color = isError ? "#f44747" : "#4ec9b0";
    setTimeout(() => el.innerText = "", 3000);
}

function openTab(tabId, btn) {
    document.querySelectorAll('.admin-tab').forEach(el => el.classList.add('hidden'));
    document.querySelectorAll('.tab-btn').forEach(el => el.classList.remove('active'));

    document.getElementById(tabId).classList.remove('hidden');
    if(btn) btn.classList.add('active');
    
    // Refresh dropdown if opening users or access tab
    if (tabId === 'tab-users' || tabId === 'tab-access'){
        updateServerDropdowns(tabId);
    }
}

async function handleLogin() {
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;
    const errorMsg = document.getElementById("error-msg");

    try {
        const res = await fetch(`${API_URL}/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password })
        });

        const data = await res.json();

        if (data.error) {
            errorMsg.innerText = data.error;
            return;
        }

        MEMORY.token = data.jwt;
        MEMORY.dek = data.token; 
        MEMORY.username = username;

        // save session to storage
        saveSession({
            token: data.jwt,
            dek: data.token,
            username: username
        });

        document.getElementById("welcome-msg").innerText = `Welcome, ${username}`;
        
        if (username.toLowerCase() === "admin") {
            document.getElementById("admin-btn").classList.remove("hidden");
        }

        loadServers(username);
        showScreen("server-screen");

    } catch (err) {
        errorMsg.innerText = "Connection Failed";
        console.error(err);
    }
}

function logout() {
    clearSession();
}

function showAdminDashboard() {
    showScreen("admin-screen");
    document.querySelector('.tab-btn').click(); 
}

function showServerList() {
    showScreen("server-screen");
    loadServers(MEMORY.username);
}

async function updateServerDropdowns(tabId) {
    let select;
    if (tabId === 'tab-users') {
        select = document.getElementById("new-user-server");
    } else if (tabId === 'tab-access') {
        select = document.getElementById("assign-server");
    }
    
    try {
        const res = await fetch(`${API_URL}/serverslist`, {
            method: "GET",
            headers: { 
                "Content-Type": "application/json",
                "Authorization": "Bearer " + MEMORY.token 
            }
        });
        const data = await res.json();
        
        const currentVal = select.value;
        
        select.innerHTML = '<option value="" disabled selected>Select Server (required)</option>';
        
        if (data.list) {
            data.list.forEach(srv => {
                const opt = document.createElement("option");
                opt.value = srv.server; 
                opt.innerText = `${srv.server} (${srv.ip})`;
                select.appendChild(opt);
            });
        }
        
        if(currentVal) select.value = currentVal;

    } catch(e) { console.error("Failed to fetch server list", e); }
}

// 1. ADD USER
async function adminAddUser() {
    const u = document.getElementById("new-user-name").value;
    const p = document.getElementById("new-user-pass").value;
    const s = document.getElementById("new-user-server").value; // Dropdown

    if(!u || !p) return setStatus("Username and Password required", true);
    if(!s) return setStatus("Server is required", true);

    try {
        const res = await fetch(`${API_URL}/register`, {
            method: "POST",
            headers: { 
                "Content-Type": "application/json",
                "Authorization": "Bearer " + MEMORY.token 
            },
            body: JSON.stringify({
                username: u,
                password: p,
                servername: s, 
                key: MEMORY.dek 
            })
        });
        const data = await res.json();
        if(data.error) return setStatus(data.error, true);
        
        setStatus("User Created Successfully!");
        document.getElementById("new-user-name").value = "";
        document.getElementById("new-user-pass").value = "";
        document.getElementById("new-user-server").value = "";
    } catch(e) { setStatus("Request Failed", true); }
}

// 2. ADD SERVER
async function adminAddServer() {
    const name = document.getElementById("srv-name").value;
    const ip = document.getElementById("srv-ip").value;
    const port = parseInt(document.getElementById("srv-port").value) || 22;
    const user = document.getElementById("srv-user").value;
    const pass = document.getElementById("srv-pass").value;

    if(!name || !ip || !user || !pass) return setStatus("All fields required", true);

    try {
        const res = await fetch(`${API_URL}/initserver`, {
            method: "POST",
            headers: { 
                "Content-Type": "application/json",
                "Authorization": "Bearer " + MEMORY.token 
            },
            body: JSON.stringify({
                servername: name,
                ip: ip,
                port: port,
                username: user,
                password: pass,
                key: MEMORY.dek 
            })
        });
        const data = await res.json();
        if(data.error) return setStatus(data.error, true);
        
        setStatus("Server Initialized!");
        document.getElementById("srv-pass").value = "";
        // Refresh dropdowns since we added a server
        updateServerDropdowns();
    } catch(e) { setStatus("Request Failed", true); }
}

// 3. ASSIGN SERVER
async function adminAssignServer() {
    const u = document.getElementById("assign-user").value;
    const s = document.getElementById("assign-server").value; // Dropdown

    if(!u || !s) return setStatus("Both User and Server required", true);

    try {
        const res = await fetch(`${API_URL}/addtouser`, {
            method: "POST",
            headers: { 
                "Content-Type": "application/json",
                "Authorization": "Bearer " + MEMORY.token 
            },
            body: JSON.stringify({
                username: u,
                servername: s
            })
        });
        const data = await res.json();
        if(data.error) return setStatus(data.error, true);
        
        setStatus(`Server ${s} assigned to ${u}`);
    } catch(e) { setStatus("Request Failed", true); }
}

async function loadServers(username) {
    const listDiv = document.getElementById("server-list");
    listDiv.innerHTML = "<p style='color:#666;'>Loading...</p>";
    
    try {
        const res = await fetch(`${API_URL}/allowedservers`, {
            method: "POST",
            headers: { 
                "Content-Type": "application/json",
                "Authorization": "Bearer " + MEMORY.token 
            },
            body: JSON.stringify({ username })
        });

        const data = await res.json();
        
        listDiv.innerHTML = "";
        if(!data.allowed || data.allowed.length === 0) {
            listDiv.innerHTML = "<p style='color:#888;'>No servers assigned.</p>";
            return;
        }

        data.allowed.forEach(serverName => {
            const btn = document.createElement("button");
            btn.className = "server-btn";
            btn.innerHTML = `<span style="color: #4ec9b0; font-weight: bold;">> ${serverName}</span> <span style="float: right; color: #666;">SSH</span>`;
            btn.onclick = () => startSession(serverName);
            listDiv.appendChild(btn);
        });

    } catch (err) {
        listDiv.innerHTML = "Failed to load servers";
    }
}

function startSession(serverName) {
    showScreen("terminal-screen");

    const term = new Terminal({
        cursorBlink: true,
        convertEol: true,
        theme: { background: '#1e1e1e' },
        fontFamily: 'monospace',
        fontSize: 14
    });
    const fitAddon = new FitAddon.FitAddon();
    term.loadAddon(fitAddon);
    
    term.open(document.getElementById("terminal-container"));
    fitAddon.fit();
    term.focus();

    MEMORY.term = term;

    const ws = new WebSocket(`${WS_URL}/ws/ssh`);
    MEMORY.socket = ws;

    ws.onopen = () => {
        term.write(`\r\n\x1b[32m[SimplePAM]\x1b[0m Establishing secure tunnel to ${serverName}...\r\n`);
        const handshake = {
            type: "handshake",
            token: MEMORY.token,
            dek: MEMORY.dek,
            target: serverName
        };
        ws.send(JSON.stringify(handshake));
    };

    ws.onmessage = (event) => term.write(event.data);
    term.onData((data) => ws.send(data));
    
    ws.onclose = () => {
        term.write("\r\n\x1b[31m[Disconnected]\x1b[0m\r\n");
    };
    
    window.onresize = () => fitAddon.fit();
}

function disconnect() {
    if (MEMORY.socket) MEMORY.socket.close();
    if (MEMORY.term) MEMORY.term.dispose();
    showScreen("server-screen");
}
