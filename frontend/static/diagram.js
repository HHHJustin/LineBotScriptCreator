export function createDiagram($) {
    if (window.myDiagram) {
        return window.myDiagram;
      }
    
      const myDiagram = $(go.Diagram, 'myDiagramDiv');
    
      window.myDiagram = myDiagram;
      return myDiagram;
}
