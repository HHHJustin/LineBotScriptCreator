export function setupDragFunctionality(myDiagram) {
    let isDragging = false;
    let startPoint = null;

    const diagramDiv = document.getElementById("myDiagramDiv");
    diagramDiv.addEventListener("mousedown", (event) => {
        if (event.target === diagramDiv) {
            isDragging = true;
            startPoint = { x: event.clientX, y: event.clientY };
        }
    });

    diagramDiv.addEventListener("mousemove", (event) => {
        if (isDragging) {
            const dx = event.clientX - startPoint.x;
            const dy = event.clientY - startPoint.y;
            myDiagram.commandHandler.moveParts(myDiagram.selection, dx, dy);
            startPoint = { x: event.clientX, y: event.clientY };
        }
    });

    diagramDiv.addEventListener("mouseup", () => {
        isDragging = false;
    });

    diagramDiv.addEventListener("mouseleave", () => {
        isDragging = false;
    });
}
