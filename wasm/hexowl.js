
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

function begin() {
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("wasm/hexowl.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
  });

  document.querySelector("#inputField").addEventListener("keydown", (e) => {
    if (e.keyCode != 13) return;
    let input = e.target.value;
    let out = hexowlPrompt(input);
    e.target.value = "";
    printOutput(">: " + input);
    printOutput(out);
  });
}

document.addEventListener("DOMContentLoaded", function() {
  if(WebAssembly) {
    begin();
  } else {
    console.log("WebAssembly is not supported in your browser")
  }
});