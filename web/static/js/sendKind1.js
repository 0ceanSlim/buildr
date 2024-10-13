document.getElementById("message-form").addEventListener("submit", async (event) => {
    event.preventDefault();
  
    const message = document.getElementById("message").value;
  
    try {
      const unsignedEvent = {
        kind: 1,
        content: message,
        created_at: Math.floor(Date.now() / 1000),
        tags: [],
      };
  
      if (!window.nostr) {
        alert("Nostr extension not available.");
        return;
      }
  
      const signedEvent = await window.nostr.signEvent(unsignedEvent);
      console.log("Signed Event:", signedEvent);
  
      const requiredFields = ["id", "pubkey", "kind", "content", "created_at", "tags", "sig"];
      const isValid = requiredFields.every((field) => signedEvent.hasOwnProperty(field) && signedEvent[field]);
  
      if (!isValid) {
        throw new Error("The signed event is missing required fields.");
      }
  
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
  
      const data = await response.json();
      console.log("Message broadcasted:", data);
  
      // Send the response data to the Kind1ResponseHandler using HTMX
      fetch("/kind1-response", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
        headers: { "HX-Request": "true" }, // Inform the server that this is an HTMX request.
      });
  
      alert("Message sent successfully.");
    } catch (error) {
      console.error("Error sending message:", error);
      alert(`Failed to send the message: ${error.message}`);
    }
  });
  