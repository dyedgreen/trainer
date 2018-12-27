"use strict";

// Api wrapper

function apiPost(endpoint, data, callback) {
  if (typeof data !== "object") throw "No data object given";

  let req = new XMLHttpRequest();
  req.open("POST", `/api${endpoint}`, true);
  let reqData = new FormData();
  Object.keys(data).forEach(key => reqData.append(key, data[key]));
  req.onload = () => {
    let res = JSON.parse(req.responseText);
    if (typeof res !== "object" || !res.hasOwnProperty("error") || !res.hasOwnProperty("value")) {
      throw "Corrupted response data";
    }
    if (typeof callback === "function") callback(res);
  };
  req.onerror = () => {
    if (typeof callback === "function") callback({error:"Server did not respond", value:null});
  };
  req.send(reqData);
}


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
ui.modal = document.getElementById("modal");
ui.buttons = {};
["run", "correct", "wrong"].forEach(btn => ui.buttons[btn] = document.getElementById(`button-${btn}`));


// App ui logic

function showModal(title, text, button, callback) {
  ui.modal.setAttribute("title", title);
  ui.modal.setAttribute("text", text);
  ui.modal.setAttribute("button", button);
  ui.modal.onclose = () => {
    callback();
    ui.modal.onclose = () => {};
  };
  ui.modal.show();
}

function showError(msg) {
  showModal("An Error Occurred", msg, "Done", () => {});
}

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
  setTimeout(() => {
    // Timeout after 30 seconds
    worker.terminate();
  }, 30000);
}
ui.buttons.run.onclick = run;

function save() {
  // Save the problem when edited
  apiPost("/problem/update", {
    id: problemId,
    title: ui.problem.title,
    question: ui.problem.question,
    solution: ui.problem.solution,
  }, res => {
    if (res.error) {
      showError(res.error);
    } else {
      // Update problem id
      problemId = +res.value;
    }
  });
}
ui.problem.onsave = save;

function submit(correct) {
  // Submit the session
  apiPost("/problem/submit", {
    id: problemId,
    code: ui.editor.getValue(),
    time: ui.timer.elapsed,
    solved: correct ? "1" : "0",
  }, res => {
    if (res.error) {
      showError(res.error);
    } else {
      showModal(
        "Session Recorded",
        correct ? "Great work today, see you tomorrow!" : "Practice makes perfect, try again tomorrow!",
        "Attempt Next Problem",
        () => window.location.href = "/app"
      );
    }
  });
}
ui.buttons.correct.onclick = () => submit(true);
ui.buttons.wrong.onclick = () => submit(false);

function loadProblem() {
  apiPost("/problem/next", {}, res => {
    if (res.error) {
      showModal(
        "Could Not Load Problem",
        res.error,
        "Try Again",
        loadProblem
      );
    } else if (res.value) {
      // A new problem is already scheduled
      problemId = +res.value.id;
      ui.problem.setTitle(res.value.title);
      ui.problem.setQuestion(res.value.question);
      ui.problem.setSolution(res.value.solution);
    } else {
      showModal(
        "No Problem Scheduled",
        "Please find a new problem and try to solve it!",
        "Done",
        () => {}
      );
    }
  });
}


// Load problem from server

let problemId = -1; // Is initialized by server, -1 <=> this will be a new problem

loadProblem();
