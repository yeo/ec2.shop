
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
  width: '100%',
  search: true,

  sort: true,

  style: {
    header: {
      'font-size': '0.8rem',
      'color': 'red',
    },
    th: {
      'font-size': '0.8rem',
      'word-wrap': 'break-word',
    },
  },

  columns: [
    {
      name: 'Type',
      width: '90px',
    },
    {
      name: 'Mem',
      width: '70px',
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

    {
      name: 'vCPUS',
      width: '60px',
    },
    {
      name: 'Storage',
      width: '70px',
    },
    {
      name: 'Network',
      width: '90px',
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
    {
      name: 'Hourly Price',
      width: '90px',
    },
    {
      name: 'Monthly',
      width: '80px',
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
      width: '80px',
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
      name: "Reserved 1y",
      width: '80px',
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
      name: "Reserved 3y",
      width: '80px',
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
      name: "1y Convertible Reser",
      width: '80px',
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
      name: "1y Convertible Reser",
      width: '80px',
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
  ],

  data: window._pricedata,
}).render(document.getElementById('price-grid'));
