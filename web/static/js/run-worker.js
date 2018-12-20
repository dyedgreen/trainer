// Worker that executes a given js string and returns
// all console logs as messages

console.log = (...args) => postMessage(args);
onmessage = msg => {
  try {
    eval("".concat(msg.data));
  } catch (e) {
    console.log(e.toString());
  }
}
