// https://github.com/ines/termynal 
var termynal = new Termynal('#termynal');

// data grid
var defaultColumns = [
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
  {
    name: 'Monthly',
    sort: {
      compare: (a, b) => {
        const code = (x) => parseFloat(x)

        if (code(a) > code(b)) {
          return 1;
        } else if (code(b) > code(a)) {
          return -1;
        } else {
          return 0;
        }
      }
    }
  }
];
var spotColumns = [
  {
    name: 'Spot Price',
    sort: {
      compare: (a, b) => {
        const code = (x) => parseFloat(x)

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
  'Spot Savings',
  'Spot Reclaim Rate',
];
var grid = new gridjs.Grid({
  search: true,
  sort: true,
  columns: defaultColumns,
  data: window._pricedata,
});

grid.render(document.getElementById('price-grid'));

document.getElementById('show-spot').onclick = function() {
  var columns = [];
  if (this.checked) {
    grid.updateConfig({
      columns: defaultColumns.concat(spotColumns),
      data: window._pricedatawithspot,
    });
  } else {
    grid.updateConfig({
      columns: defaultColumns,
      data: window._pricedata,
    });
  }
  grid.forceRender();
};
