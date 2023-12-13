
var inputHistory = [];
var inputHistoryIndex = 0;

function printOutput(out) {
  let outDiv = document.querySelector("#output");
  let lines = out.split("\n");

  lines.forEach(line => {
    let newln = document.createElement("pre");

    if (line.length > 0)
      newln.innerText = line;
    else
      newln.innerText = " ";

    outDiv.append(newln);
    newln.scrollIntoView({ behavior: "smooth", block: "end", inline: "nearest" });
  });
}

function loadEnv() {

}

function saveEnv(name = "") {

}

function begin() {
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("wasm/hexowl.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    let out = hexowlPrompt("help");
    printOutput(">: help");
    printOutput(out);
  });

  document.querySelector("#inputField").addEventListener("keydown", (e) => {
    switch (e.key) {
      case "Enter": {
        let input = e.target.value;
        let out = hexowlPrompt(input);
        inputHistory.unshift(input);
        inputHistoryIndex = -1;
        e.target.value = "";
        printOutput(">: " + input);
        printOutput(out);
        break;
      }
      case "ArrowUp": {
        if (inputHistoryIndex < inputHistory.length - 1) {
          inputHistoryIndex += 1;
          e.target.value = inputHistory[inputHistoryIndex];
          e.target.setSelectionRange(e.target.value.length, e.target.value.length);
        }
        e.preventDefault();
        break;
      }
      case "ArrowDown": {
        if (inputHistoryIndex > 0) {
          inputHistoryIndex -= 1;
          e.target.value = inputHistory[inputHistoryIndex];
          e.target.setSelectionRange(e.target.value.length, e.target.value.length);
        } else {
          inputHistoryIndex = -1;
          e.target.value = "";
        }
        e.preventDefault();
        break;
      }
    }
  });
}

document.addEventListener("DOMContentLoaded", function() {
  if(WebAssembly) {
    begin();
  } else {
    console.log("WebAssembly is not supported in your browser")
  }
});