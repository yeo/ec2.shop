// https://github.com/ines/termynal 
var termynal = new Termynal('#termynal');

function weightNetwork(value) {
  if (value == 'NA') {
    return -1001
  }

  if (value == 'Very Low') {
    return -1000
  }

  if (value == 'Low') {
    return -998
  }

  if (value == 'Low to Moderate') {
    return -997
  }

  if (value == 'Moderate') {
    return -996
  }

  if (value == 'High') {
    return -995
  }

  if (value.startsWith('Up to')) {
    const m = value.match(/\d+/)
    if (m) {
      return parseInt(m[0]) - 0.5
    }
  }

  const m = value.match(/\d+/)
  if (m) {
    return parseInt(m[0])
  }

  return -9000
}

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
    {
      name: 'Network',
      sort: {
        compare: (a, b) => {
          const wa = weightNetwork(a)
          const wb = weightNetwork(b)
          if (wa > wb) {
            return 1
          }

          if (wa < wb) {
            return -1
          }

          return 0
        }
      }
    },
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
    },
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
    }
  ],

  data: window._pricedata,
}).render(document.getElementById('price-grid'));
