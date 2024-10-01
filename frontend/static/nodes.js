import { createDiagram } from './diagram.js';
import { setupNodeTemplate } from './templates.js';
import { fetchData } from './fetchData.js';

async function init() {
    const $ = go.GraphObject.make;
    const myDiagram = createDiagram($);
    if (!(myDiagram instanceof go.Diagram)) {
        console.error("Failed to create GoJS Diagram");
        return;
    }
    setupNodeTemplate(myDiagram, $);
    await fetchData(myDiagram);

    myDiagram.addDiagramListener("SelectionMoved", async function(e) {
        e.subject.each(async function(part) {
            if (part instanceof go.Node) {
                const data = {
                    currentNodeID: part.data.key,
                    locX: part.location.x,
                    locY: part.location.y,
                };
                try {
                    const response = await fetch('/nodes/updatelocation', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify(data),
                    });

                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    console.log('Node location updated successfully');
                } catch (error) {
                    console.error('Error updating node location:', error);
                }
            }
        });
    });
}

window.addEventListener('DOMContentLoaded', init);
