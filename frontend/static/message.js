import { getCurrentIDAndTypeFromURL } from './utils.js';

export function makeEditable(td, messageID) {
    const originalContent = td.innerText;
    const input = document.createElement('input');
    input.type = 'text';
    input.value = originalContent;
    input.onkeydown = async function(event) {
        if (event.key === 'Enter') {
            const messageContent = input.value;
            const response = await fetch('/messages/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    messageID: messageID,
                    messageContent: messageContent,
                }),
            });
            if (response.ok) {
                const jsonResponse = await response.json();
                console.log('Message updated:', jsonResponse);
                td.innerText = messageContent; 
            } else {
                console.error('Error updating message');
            }
        }
    };
    td.innerHTML = '';
    td.appendChild(input);
    input.focus();
}

export function deleteMessage(event, messageID) {
    event.preventDefault();  
    const url = window.location.pathname;  
    const segments = url.split('/');  
    const currentID = segments[3]; 
    if (!currentID) {
        console.error("Node ID is not found in the URL.");
        return;
    }
    const currentIDInt = parseInt(currentID, 10);
    if (isNaN(currentIDInt)) {
        console.error("Node ID is not a valid number.");
        return;
    }
    const messageRow = event.target.closest('tr'); 
    const indexCell = messageRow.querySelector('td:first-child'); 
    const messageIndex = indexCell.textContent; 
    const messageIndexInt = parseInt(messageIndex, 10); 
    const data = {
        messageID: messageID,
        currentNodeID: currentIDInt,
        messageIndex: messageIndexInt
    };
    console.log(data)
    fetch('/messages/delete', {
        method: 'POST',  
        headers: {
            'Content-Type': 'application/json',  
        },
        body: JSON.stringify(data)  
    })
    .then(response => {
        if (!response.ok) {
            throw new Error("Network response was not ok");
        }
        return response.json();  
    })
    .then(data => {
        console.log('Success:', data);
        window.location.reload();  
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

export  function updatePlaceholder() {
    const messageType = document.getElementById("messageTypeSelect").value;
    const messageInput = document.getElementById("messageContentInput");
    if (messageType === "Text") {
        messageInput.placeholder = "Enter Text message";
    } else if (messageType === "Image") {
        messageInput.placeholder = "Enter Image URL";
    } else if (messageType === "FlexMessage") {
        messageInput.placeholder = "Enter FlexMessage JSON";
    } else {
        messageInput.placeholder = "Please Choose Image Type";
    }
}

document.getElementById("messageContentInput").addEventListener("keydown", function(event) {
    if (event.key === "Enter") {
        event.preventDefault();  
        submitMessage();  
    }
});

export async function submitMessage() {
    const { currentID } = getCurrentIDAndTypeFromURL(); 
    const currentIDInt = parseInt(currentID, 10);
    var messageType = document.getElementById("messageTypeSelect").value;
    var messageContent = document.getElementById("messageContentInput").value;

    if (!messageType || !messageContent) {
        alert("請填寫完整的資料");
        return;
    }

    var data = {
        currentNodeID: currentIDInt,
        messageType: messageType,
        messageContent: messageContent
    };

    try {
        const response = await fetch('/messages/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const result = await response.json();
        console.log('Success:', result);
        window.location.reload();

    } catch (error) {
        console.error('Error:', error);
    }
}
