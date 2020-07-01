// https://github.com/ines/termynal 
var termynal = new Termynal('#termynal');

// data grid
new gridjs.Grid({
  search: true,
  sort: true,

  columns: [
    'Instance Type',
    {
      name: 'Memory',
      sort: {
        compare: (a, b) => {
          const code = (x) => parseFloat(x.split(' ')[0])

          if (code(a) > code(b)) {
            return 1;
          } else if (code(b) > code(a)) {
            return -1;
          } else {
            return 0;
          }
        }
      }
    },

    'vCPUS',
    'Storage',
    'Network',
    'Hourly Price',
    'Monthly',
  ],

  data: window._pricedata,
}).render(document.getElementById('price-grid'));
