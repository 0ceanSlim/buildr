document.getElementById("message-form").addEventListener("submit", async (event) => {
    event.preventDefault();
  
    const message = document.getElementById("message").value;
    const statusElement = document.getElementById("message-status");
  
    // Function to update status
    function updateStatus(text) {
        statusElement.textContent = text;
        console.log("Status updated:", text); // For debugging
    }

    try {
        updateStatus("Preparing message...");

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
        const signedEvent = await window.nostr.signEvent(unsignedEvent);
        console.log("Signed Event:", signedEvent);
  
        updateStatus("Sending to relays...");
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
        const relayResults = await response.json();
        console.log("Message broadcasted:", relayResults);
  
        // Update the status element with the relay results
        updateStatusWithRelayResults(statusElement, relayResults);
  
    } catch (error) {
        console.error("Error sending message:", error);
        updateStatus(`Error: ${error.message}`);
    }
});

function updateStatusWithRelayResults(statusElement, relayResults) {
    let resultHtml = "<h3>Relay Results:</h3><ul>";
    
    for (const [url, status] of Object.entries(relayResults)) {
        resultHtml += `<li>${url}: ${status}</li>`;
    }
    
    resultHtml += "</ul>";
    statusElement.innerHTML = resultHtml;
    console.log("Status updated with relay results"); // For debugging
}