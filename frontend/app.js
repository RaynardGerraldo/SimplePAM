const MEMORY = {
    token: null,  // jwt
    dek: null,    // dek
    socket: null, 
    term: null
};

const API_URL = "http://localhost:8080";
const WS_URL = "ws://localhost:8080";

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

        document.getElementById("login-screen").classList.add("hidden");
        document.getElementById("server-screen").classList.remove("hidden");
        
        loadServers(username);

    } catch (err) {
        errorMsg.innerText = "Connection Failed";
        console.error(err);
    }
}

// get server list
async function loadServers(username) {
    const listDiv = document.getElementById("server-list");
    
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
        
        // button for server list
        listDiv.innerHTML = "";
        data.allowed.forEach(serverName => {
            const btn = document.createElement("button");
            btn.className = "server-btn";
            btn.innerText = `Connect to ${serverName}`;
            btn.onclick = () => startSession(serverName);
            listDiv.appendChild(btn);
        });

    } catch (err) {
        alert("Failed to load servers");
    }
}


function startSession(serverName) {
    // ui switch
    document.getElementById("server-screen").classList.add("hidden");
    document.getElementById("terminal-screen").classList.remove("hidden");

    // init terminal
    const term = new Terminal({
        cursorBlink: true,
        theme: { background: '#1e1e1e' }
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
        term.write("\r\nConnecting to " + serverName + "...\r\n");

        const handshake = {
            type: "handshake",
            token: MEMORY.token,
            dek: MEMORY.dek,
            target: serverName
        };
        ws.send(JSON.stringify(handshake));
    };

    ws.onmessage = (event) => {
        term.write(event.data);
    };

    term.onData((data) => {
        ws.send(data);
    });

    ws.onclose = () => {
        term.write("\r\n\r\n--- Connection Closed ---\r\n");
    };
    
    window.onresize = () => fitAddon.fit();
}

function disconnect() {
    if (MEMORY.socket) {
        MEMORY.socket.close();
    }
    if (MEMORY.term) {
        MEMORY.term.dispose();
    }
    document.getElementById("terminal-screen").classList.add("hidden");
    document.getElementById("server-screen").classList.remove("hidden");
}
