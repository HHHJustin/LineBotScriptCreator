import { createDiagram } from './diagram.js';
import { setupNodeTemplate } from './templates.js';
import { fetchData } from './fetchData.js';

async function init() {
    const $ = go.GraphObject.make;
    const myDiagram = createDiagram($);
    setupNodeTemplate(myDiagram, $);
    await fetchData(myDiagram);
}

window.addEventListener('DOMContentLoaded', init);
