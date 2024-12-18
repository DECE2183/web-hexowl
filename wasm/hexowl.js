
var inputHistory = [];
var inputHistoryIndex = 0;

function clearOutput() {
  let outDiv = document.querySelector("#output");
  outDiv.innerHTML = "";
}

async function openFile() {
  let p = new Promise((resolve) => {
    let input = document.createElement("input");
    input.type = "file";
    input.onchange = (e) => {
      let file = e.target.files[0];
      if (file == null) {
        resolve(null);
      }

      let fileReader = new FileReader();

      fileReader.onloadend = () => {
        resolve(fileReader.result);
      }

      fileReader.readAsText(file);
    };
    input.click();
  });

  return await p;
}

function saveFile(name, content) {
  let content_blob = new Blob([content], { type: "text/plain" });
  let url = window.URL.createObjectURL(content_blob);

  let a = document.createElement("a");
  a.style = "display: none";
  document.body.appendChild(a);
  a.href = url;
  a.download = `${name}.json`;
  a.click();
  window.URL.revokeObjectURL(url);
  a.remove();
}

function highlightLn(ln, text) {
  let escBegin = false;
  let escSpan = document.createElement("span");
  let escArgs = "";

  if (text.length == 0) {
    return;
  }

  for (let ch of text) {
    if (ch == '\x1B') {
      escBegin = true;
      ln.append(escSpan);
      escSpan = document.createElement("span");
      continue;
    }
    
    if (escBegin) {
      switch (ch) {
      case '[':
        escArgs = "";
        continue;
      case 'm':
        escSpan.classList.add("esc-"+escArgs.replaceAll(';', '-'));
        escBegin = false;
        continue;
      default:
        escArgs += ch;
        continue;
      }
    } else {
      escSpan.innerText += ch;
    }
  }

  ln.append(escSpan);
}

function printOutput(out) {
  let outDiv = document.querySelector("#output");
  let lines = out.split("\n");

  lines.forEach(line => {
    let newln = document.createElement("pre");
    highlightLn(newln, line);
    outDiv.append(newln);
    newln.scrollIntoView({ behavior: "smooth", block: "end", inline: "nearest" });
  });
}

function begin() {
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("wasm/hexowl.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    hexowlPrompt("help");
  });

  document.querySelector("#inputField").addEventListener("keydown", async (e) => {
    switch (e.key) {
      case "Enter": {
        let input = e.target.value;
        if (input.length == 0) break;
        hexowlPrompt(input);
        inputHistory.unshift(input);
        inputHistoryIndex = -1;
        e.target.value = "";
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

  document.querySelector("#inputEnterBtn").addEventListener("click", (e) => {
    let inputField = document.querySelector("#inputField");
    if (inputField.value.length == 0) return;
    hexowlPrompt(inputField.value);
    inputHistory.unshift(inputField.value);
    inputHistoryIndex = -1;
    inputField.value = "";
    inputField.focus();
  });
}

document.addEventListener("DOMContentLoaded", function() {
  if(WebAssembly) {
    begin();
  } else {
    console.log("WebAssembly is not supported in your browser")
  }
});