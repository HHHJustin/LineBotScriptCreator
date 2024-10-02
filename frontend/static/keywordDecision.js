import { getCurrentIDAndTypeFromURL } from './utils.js';

export async function submitDeleteKWDecision(event, keywordDecisionID) {
    event.preventDefault();  
    const { currentID } = getCurrentIDAndTypeFromURL(); 
    const currentIDInt = parseInt(currentID, 10);

    const data = {
        keywordDecisionID: keywordDecisionID,
        currentNodeID: currentIDInt
    };

    try {
        const response = await fetch('/keywordDecisions/delete', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',  
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            throw new Error("Network response was not ok");
        }

        const result = await response.json();
        console.log('Success:', result);
        window.location.reload();  

    } catch (error) {
        console.error('Error deleting message:', error);
        alert('Error deleting message: ' + error.message); 
    }
}

export function makeKeywordEditable(td, keywordDecisionID) {
    const originalKeyword = td.innerText; 
    const input = document.createElement('input'); 
    input.type = 'text';
    input.value = originalKeyword;

    input.onkeydown = async function(event) {
        if (event.key === 'Enter') {
            const newKeyword = input.value;
            try {
                const response = await fetch('/keywordDecisions/update', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        KeywordDecisionID: keywordDecisionID,
                        Keyword: newKeyword,
                    }),
                });

                if (response.ok) {
                    const jsonResponse = await response.json();
                    console.log('Keyword updated successfully:', jsonResponse);
                    td.innerText = newKeyword; 
                } else {
                    console.error('Error updating keyword');
                    alert('Error updating keyword');
                }
            } catch (error) {
                console.error('Error:', error);
                alert('Error updating keyword: ' + error.message);
            }
        }
    };
    td.innerHTML = ''; 
    td.appendChild(input);
    input.focus(); 
}