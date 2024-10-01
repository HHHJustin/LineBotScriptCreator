export function createDiagram($) {
    if (window.myDiagram) {
        return window.myDiagram;
      }
      const myDiagram = $(go.Diagram, "myDiagramDiv", {
        "undoManager.isEnabled": true, 
    });
    
      window.myDiagram = myDiagram;
      return myDiagram;
}
