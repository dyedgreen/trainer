"use strict";

// Custom HTML Elements

class ModalWidget extends HTMLElement {
  constructor() {
    super();
    // State
    this.innerHTML = '<div><h1></h1><p></p><a class="button"></a></div>';
    this.titleElem = this.childNodes[0].childNodes[0];
    this.textElem = this.childNodes[0].childNodes[1];
    this.buttonElem = this.childNodes[0].childNodes[2];
    // Actions
    this.buttonElem.onclick = () => {this.close()};
    // Draw modal
    this.titleElem.innerHTML = this.hasAttribute("title") ? this.getAttribute("title") : "";
    this.textElem.innerHTML = this.hasAttribute("text") ? this.getAttribute("text") : "";
    this.buttonElem.innerHTML = this.hasAttribute("button") ? this.getAttribute("button") : "";
  }

  setAttribute(...args) {
    HTMLElement.prototype.setAttribute.call(this, ...args);
    if (["title", "text", "button"].indexOf(args[0]) !== -1) {
      this[`${args[0]}Elem`].innerHTML = this.hasAttribute(args[0]) ? this.getAttribute(args[0]) : "";
    }
  }

  show() {
    this.classList.remove("hidden");
  }

  close() {
    this.classList.add("hidden");
    this.onclose();
  }

  onclose() {
    // Overwrite this function to register a handler
  }
}
customElements.define("modal-widget", ModalWidget);

class TimerWidget extends HTMLElement {
  constructor() {
    super();
    // State
    this.innerHTML = '<h1 class="font-mono"></h1><a class="icon-play"></a>';
    this.display = this.childNodes[0];
    this.button = this.childNodes[1];
    this.paused = true;
    this.locked = false;
    this.elapsed = 0;
    this.timeout = null;
    // Actions
    this.button.onclick = () => {this.toggle()};
    // Initial draw
    this.display.innerHTML = this.displayTime;
  }

  toggle() {
    if (this.locked) this.paused = false;
    this.paused = !this.paused;
    clearTimeout(this.timeout);
    if (this.paused) {
      this.button.classList.add("icon-play");
      this.button.classList.remove("icon-pause");
    } else {
      this.button.classList.remove("icon-play");
      this.button.classList.add("icon-pause");
      setTimeout(() => {this.tick()}, 1000);
    }
    this.ontoggle(this.paused);
  }

  get duration() {
    return this.hasAttribute("duration") ? +this.getAttribute("duration") : 30*60;
  }

  get remaining() {
    return this.duration - this.elapsed;
  }

  get displayTime() {
    const minutes = Math.floor(Math.abs(this.remaining) / 60);
    const seconds = Math.abs(this.remaining) % 60;
    const neg = this.remaining < 0;
    return (neg ? "-" : "").concat(minutes) + (seconds < 10 ? ":0" : ":").concat(seconds);
  }

  tick() {
    if (this.paused || this.locked) return;
    this.elapsed ++;
    this.display.innerHTML = this.displayTime;
    if (this.elapsed === this.duration) {
      this.onfinish();
    }
    this.timeout = setTimeout(() => {this.tick()}, 1000);
  }

  reset() {
    this.elapsed = 0;
    this.display.innerHTML = this.displayTime;
    this.paused = false;
    this.toggle();
  }

  lock() {
    this.classList.add("locked");
    this.locked = true;
    this.toggle();
  }

  ontoggle(paused) {
    // Overwrite this function to register a handler
  }

  onfinish() {
    // Overwrite this function to register a handler
    this.reset();
  }
}
customElements.define("timer-widget", TimerWidget);

class ProblemWidget extends HTMLElement {
  constructor() {
    super();
    // State
    this.innerHTML = '<h1></h1><p></p><h2 class="hidden">Solution</h2><p class="hidden"></p><a class="button icon-done">Show Solution</a><a class="button icon-edit">Edit</a>';
    this.titleElem = this.childNodes[0];
    this.questionElem = this.childNodes[1];
    this.solutionHead = this.childNodes[2];
    this.solutionElem = this.childNodes[3];
    this.buttonDone = this.childNodes[4];
    this.buttonEdit = this.childNodes[5];
    this.edit = false;
    // Actions
    this.buttonDone.onclick = () => {this.setDone()}
    this.buttonEdit.onclick = () => {this.toggleEdit()};
    // Initial draw
    this.titleElem.innerHTML = this.hasAttribute("title") ? this.getAttribute("title") : "";
    this.questionElem.innerHTML = this.hasAttribute("question") ? this.getAttribute("question") : "";
    this.solutionElem.innerHTML = this.hasAttribute("solution") ? this.getAttribute("solution") : "";
  }

  setTitle(str) {
    this.titleElem.innerHTML = "".concat(str);
  }

  setQuestion(str) {
    this.questionElem.innerHTML = "".concat(str);
  }

  setSolution(str) {
    this.solutionElem.innerHTML = "".concat(str);
  }

  get title() {
    return "".concat(this.titleElem.innerHTML);
  }

  get question() {
    return "".concat(this.questionElem.innerHTML);
  }

  get solution() {
    return "".concat(this.solutionElem.innerHTML);
  }

  sanitize () {
    [this.titleElem, this.questionElem, this.solutionElem].forEach(elem => {
      elem.innerHTML = "".concat(elem.innerHTML).replace(
        /\<[^>/]+\>/g, "\n").replace(
        /\<[^>]+\>/g, "");
    });
  }

  toggleEdit() {
    this.edit = !this.edit;
    let contenteditable = "false";
    if (this.edit) {
      contenteditable = "true";
      this.buttonEdit.classList.add("icon-save");
      this.buttonEdit.classList.remove("icon-edit");
    } else {
      this.buttonEdit.classList.add("icon-edit");
      this.buttonEdit.classList.remove("icon-save");
      this.sanitize();
      this.onsave();
    }
    this.buttonEdit.innerHTML = this.edit ? "Save" : "Edit";
    this.titleElem.setAttribute("contenteditable", contenteditable);
    this.questionElem.setAttribute("contenteditable", contenteditable);
    this.solutionElem.setAttribute("contenteditable", contenteditable);
  }

  setDone() {
    this.solutionHead.classList.remove("hidden");
    this.solutionElem.classList.remove("hidden");
    this.buttonDone.classList.add("hidden");
    this.ondone();
  }

  onsave() {
    // Overwrite this function to register a handler
  }

  ondone() {
    // Overwrite this function to register a handler
  }
}
customElements.define("problem-widget", ProblemWidget);

class ConsoleWidget extends HTMLElement {
  constructor() {
    super();
    this.messages = [];
    this.render();
  }

  render() {
    this.innerHTML = "> ".concat(this.messages.join("<br>&gt; "));
    this.scrollTop = this.scrollHeight - this.clientHeight;
    if (this.messages.length === 0) this.innerHTML = "> no messages";
  }

  log(...msg) {
    this.messages.push(msg.reduce((str, obj) => {
      return str + " ".concat(typeof obj === "string" ? obj : JSON.stringify(obj))
    }, ""));
    this.render();
  }

  clear() {
    this.messages = [];
    this.render();
  }
}
customElements.define("console-widget", ConsoleWidget);
