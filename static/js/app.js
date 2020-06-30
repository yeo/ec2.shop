// https://github.com/ines/termynal 
var termynal = new Termynal('#termynal');

// data grid
new gridjs.Grid({
  search: true,
  sort: true,

  columns: ['Instance Type', 'Memory', 'vCPUS', 'Storage', 'Network', 'Price'],

  data: window._pricedata,
}).render(document.getElementById('price-grid'));
