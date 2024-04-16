function median(a) {
    return a.reduce((acc, i) => acc+i) / (a.length * 1000_000)
  }
  
  function get90th(a) {
    return a.sort((a, b) => a - b).slice(a.length * 0.9)
  }
  
  function get99th(a) {
    return a.sort((a, b) => a - b).slice(a.length * 0.99)
  }
  
  function absIt(a) {
    return a.map(x => Math.abs(x))
  }
  
  C = [
    
  ]
  
  console.log("50th: ", median(absIt(C)))
  console.log("90th: ", median(get90th(absIt(C))))
  console.log("99th: ", median(get99th(absIt(C))))
  
  