"use strict";

// Initialize app, gather all components

let ui = {};

ui.editor = CodeMirror.fromTextArea(
  document.getElementById("code"),
  {
    lineNumbers: true,
    smartIndent: true,
    tabSize: 2,
    readOnly: true,
    theme: "material",
    mode: "javascript",
});
ui.editor.setSize("100%", "100%");
ui.timer = document.getElementById("timer");
ui.problem = document.getElementById("problem");
ui.console = document.getElementById("console");
ui.buttons = {};
["run", "correct", "wrong"].forEach(btn => ui.buttons[btn] = document.getElementById(`button-${btn}`));


// App ui logic

function finish() {
  // Called on timer stop and view-solution
  ui.editor.setOption("readOnly", true);
  ui.timer.lock();
}
ui.timer.onfinish = finish;
ui.problem.ondone = finish;

function toggle(paused) {
  // Updates writable state of editor
  ui.editor.setOption("readOnly", paused);
}
ui.timer.ontoggle = toggle;


function run() {
  // Executes the current editor code
  ui.console.clear();
  let worker = new Worker("/js/run-worker.js");
  worker.onerror = () => ui.console.log("An error occurred!");
  worker.onmessage = msg => ui.console.log(...msg.data);
  worker.postMessage(ui.editor.getValue());
}
ui.buttons.run.onclick = run;

function submit(correct) {
  // Post result to server
}
ui.buttons.correct.onclick = () => submit(true);
ui.buttons.wrong.onclick = () => submit(false);


// Load problem from server
