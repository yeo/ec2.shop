new gridjs.Grid({
  search: true,

  columns: ['Instance Type', 'Memory', 'vCPUS', 'Storage', 'Network', 'Price'],
  sort: true,
  data: window._pricedata,
}).render(document.getElementById('price-grid'));
