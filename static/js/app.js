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

function compareNetwork(a, b) {
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

function compareFloat(a, b) {
  const code = (x) => parseFloat(x)

  if (code(a) > code(b)) {
    return 1;
  } else if (code(b) > code(a)) {
    return -1;
  } else {
    return 0;
  }
}

function compareRange(a, b) {
  const range = {
      "NA": -1,
      "<5%": 0,
      "5-10%": 1,
      "10-15%": 2,
      "15-20%": 3,
      ">20%": 4,
  }

  const code = (x) => range[x]

  if (code(a) > code(b)) {
    return 1;
  } else if (code(b) > code(a)) {
    return -1;
  } else {
    return 0;
  }
}

// compare two string by the first float part
function compareFloatFirst(a, b) {
  if (a == "NA") {
    return -1;
  }
  if (b == "NA") {
    return 1;
  }

  const code = (x) => parseFloat(x.split(' ')[0])

  if (code(a) > code(b)) {
    return 1;
  } else if (code(b) > code(a)) {
    return -1;
  } else {
    return 0;
  }
}

let awsSvc = "ec2";
if (document.currentScript.hasAttribute('data-svc')) {
    // svc has this format provider:svc such as aws:rds aws:rds-pg or gcp:vm
    const svc = document.currentScript.getAttribute('data-svc');
    awsSvc = svc.split(":")[1];
}
let params = new URL(document.location.toString()).searchParams;

const dataGridOptions = {}
dataGridOptions.ec2 = {
  server: {
    url: '/?json&' + params.toString(),
    then: (data) => {
      //window.history.pushState(params, 'unused', '?');
      return data.Prices.map(price => [
        price.InstanceType,
        price.Memory,
        price.VCPUS,
        price.Storage,
        price.Network,
        price.Cost,
        price.MonthlyPrice,
        price.SpotPrice,
        price.SpotReclaimRate,
        price.SpotSavingRate,
        price.Reserved1yPrice,
        price.Reserved3yPrice,
        price.Reserved1yConveritblePrice,
        price.Reserved3yConveritblePrice,
      ])
    }
  },

  columns: [
    {
      name: 'Type',
      width: '90px',
    },
    {
      name: "Mem (GiB)",
      width: '70px',
      sort: {
        compare: compareFloatFirst,
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
        compare: compareNetwork
      }
    },
    {
      name: 'Price',
      width: '90px',
      columns: [
        {
          name: 'Hourly',
          sort: {
            compare: compareFloat
          }
        },
        {
          name: 'Monthly',
          sort: {
            compare: compareFloat
          }
        },
      ],
    },
    {
      name: 'Spot',
      width: '120px',
      columns: [
        {
          name: 'Price',
          width: '30px',
          sort: {
            compare: compareFloat,
          }
        },
        {
          name: 'Reclaim',
          width: '50px',
          sort: {
            compare: compareRange,
          }
        },
        {
          name: 'Saving',
          width: '40px',
          sort: {
          compare: compareFloat
          }
        }
      ],
    },
    {
      name: "Reserved",
      width: '80px',
      columns: [{
        name: '1 year',
        width: '40px',
        sort: {
          compare: compareFloat
        }
      }, {
        name: '3 year',
        width: '40px',
        sort: {
          compare: compareFloat
        }
      }],
    },
    {
      name: "Reserved Convertible",
      width: '100px',
      columns: [{
        name: "1y",
        width: '50px',
        sort: {
          compare: compareFloat
        }
      }, {
        name: "3y",
        width: '50px',
        sort: {
          compare: compareFloat
        }
      }]
    },
  ],
}

dataGridOptions.rds = {
  server: {
    url: '/rds?json&' + params.toString(),
    then: (data) => {
      //window.history.pushState(params, 'unused', '?');
      return data.Prices.map(price => [
        price.InstanceType,
        price.Memory,
        price.VCPUS,
        price.Network,
        price.Cost,
        price.MonthlyPrice,
        price.MultiAZ,
        price.MultiAZ2,
        price.Reserved1yPrice,
        price.Reserved3yPrice,
        price.Reserved1yMultiAZ,
        price.Reserved3yMultiAZ,
      ])
    }
  },

  columns: [
    {
      name: 'Type',
      width: '90px',
    },
    {
      name: "Mem (GiB)",
      width: '70px',
      sort: {
        compare: compareFloatFirst
      }
    },

    {
      name: 'vCPUS',
      width: '60px',
    },
    {
      name: 'Network',
      width: '90px',
      sort: {
        compare: compareNetwork
      }
    },
    {
      name: 'SingleAZ',
      width: '80px',
      columns: [{
        name: 'Hourly',
        width: '40px',
        sort: {
          compare: compareFloat
        }
      }, {
        name: 'Monthly',
        width: '40px',
        sort: {
          compare: compareFloat
        }
      }],
    },
    {
      name: 'MultiAZ',
      width: '80px',
      columns: [{
        name: "1 standby",
        width: '40px',
        sort: {
          compare: compareFloat
        }
      }, {
        name: "2 standby",
        width: '40px',
        sort: {
          compare: compareFloat
        }
      }]
    },
    {
      name: "SingleAZ Reserved",
      width: '80px',
      columns: [{
        name: '1y NoUpfront',
        width: '40px',
        sort: {
          compare: compareFloat
        }
      }, {
        name: '3y Partial',
        width: '40px',
        sort: {
          compare: compareFloat
        }
      }]
    },
    {
      name: "MultiAZ Reserved",
      width: '80px',
      columns: [{
        name: '1y NoUpfront',
        width: '40px',
        sort: {
          compare: compareFloat
        }
      }, {
        name: '3y Partial',
        width: '40px',
        sort: {
          compare: compareFloat
        }

      }]
    },

  ],
}


new gridjs.Grid({
  ...(dataGridOptions[awsSvc]),
  width: '100%',
  fixedHeader: true,
  height: '800px',
  className: {
    td: 'align-top mt-2',
  },

  sort: true,

  search: {
    server: {
      url: (prev, keyword) => {
        let params = new URL(document.location.toString()).searchParams;
        params.set("filter", keyword);
        if (!params.get("region")) {
            // default to us-east-1
            // TODO: load from cookie ?
            params.set("region", "us-east-1");
        }
        //window.history.pushState(params, 'unused', '?');
        return `?json&${params.toString()}`
      }
    }
  },

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
}).render(document.getElementById('price-grid'))
