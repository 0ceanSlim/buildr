document.getElementById("message-form").addEventListener("submit", async (event) => {
    event.preventDefault();
  
    const message = document.getElementById("message").value;
    const formElement = document.getElementById("message-form");
    const statusElement = document.getElementById("message-status");
    const statusContent = document.getElementById("status-content");
  
    function updateStatus(text) {
        statusContent.textContent = text;
        console.log("Status updated:", text);
    }

    function transition(from, to) {
        from.classList.add('fade-out');
        setTimeout(() => {
            from.classList.add('hidden');
            from.classList.remove('fade-out');
            to.classList.remove('hidden');
            setTimeout(() => to.classList.add('fade-in'), 50);
        }, 300);
    }

    transition(formElement, statusElement);

    try {
        updateStatus("Preparing message...");
        await new Promise(resolve => setTimeout(resolve, 1000)); // Simulating delay

        const unsignedEvent = {
            kind: 1,
            content: message,
            created_at: Math.floor(Date.now() / 1000),
            tags: [],
        };
  
        if (!window.nostr) {
            updateStatus("Nostr extension not available.");
            return;
        }
  
        updateStatus("Signing event...");
        await new Promise(resolve => setTimeout(resolve, 1000)); // Simulating delay
        const signedEvent = await window.nostr.signEvent(unsignedEvent);
        console.log("Signed Event:", signedEvent);
  
        updateStatus("Sending to relays...");
        await new Promise(resolve => setTimeout(resolve, 1000)); // Simulating delay
        const response = await fetch("/send-signed-kind1", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(signedEvent),
        });
  
        if (!response.ok) {
            const errorMessage = await response.text();
            throw new Error(`Failed to broadcast event: ${errorMessage}`);
        }
  
        updateStatus("Processing server response...");
        await new Promise(resolve => setTimeout(resolve, 1000)); // Simulating delay
        const relayResults = await response.json();
        console.log("Message broadcasted:", relayResults);
  
        updateStatusWithRelayResults(statusContent, relayResults);
  
    } catch (error) {
        console.error("Error sending message:", error);
        updateStatus(`Error: ${error.message}`);
    }
});

document.getElementById("close-status").addEventListener("click", () => {
    const formElement = document.getElementById("message-form");
    const statusElement = document.getElementById("message-status");
    transition(statusElement, formElement);
    document.getElementById("message").value = ''; // Clear the input
});

function updateStatusWithRelayResults(statusContent, relayResults) {
    let resultHtml = "<h3>Relay Results:</h3><ul>";
    
    for (const [url, status] of Object.entries(relayResults)) {
        const emoji = status === 'Success' ? '✅' : '❌';
        resultHtml += `<li>${emoji} ${url}: ${status}</li>`;
    }
    
    resultHtml += "</ul>";
    statusContent.innerHTML = resultHtml;
    console.log("Status updated with relay results");
}

function transition(from, to) {
    from.classList.add('fade-out');
    setTimeout(() => {
        from.classList.add('hidden');
        from.classList.remove('fade-out');
        to.classList.remove('hidden');
        setTimeout(() => to.classList.add('fade-in'), 50);
    }, 300);
}