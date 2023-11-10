/**
 * @param {string} str
 * @returns {HTMLElement}
 */
function $id(str) {
  const el = document.getElementById(str);
  if (el) return el;
  throw `#${str} does not exist.`;
}
/**
 * @param {string} str
 * @returns {HTMLButtonElement}
 */
function $button(str) {
  const el = $id(str);
  if (el instanceof HTMLButtonElement) return el;
  throw `#${str} not HTMLButtonElement.`;
}
/**
 * @param {string} str
 * @returns {HTMLDivElement}
 */
function $div(str) {
  const el = $id(str);
  if (el instanceof HTMLDivElement) return el;
  throw `#${str} not HTMLDivElement.`;
}
/**
 * @param {string} str
 * @returns {HTMLInputElement}
 */
function $input(str) {
  const el = $id(str);
  if (el instanceof HTMLInputElement) return el;
  throw `#${str} not HTMLInputElement.`;
}

const bClose = $button('close');
const bLogin = $button('login');
const bLogout = $button('logout');
const bOpen = $button('open');
const bSend = $button('send');
const divOutput = $div('output');
const iInput = $input('input');

/** @type {WebSocket|null} */
let ws = null;

/**
 * @param {*} message
 */
function out(...message) {
  console.log(message);
  const d = document.createElement('div');
  d.textContent = message;
  divOutput.appendChild(d);
  divOutput.scroll(0, divOutput.scrollHeight);
}

bLogin.addEventListener('click', async function () {
  const r = await fetch('http://localhost:8080/login', {
    credentials: 'include',
  });
  r.status === 200 ? out('LOGIN SUCCESS') : out('LOGIN FAILURE');
});

bLogout.addEventListener('click', async function () {
  const r = await fetch('http://localhost:8080/logout', {
    credentials: 'include',
  });
  r.status === 200 ? out('LOGOUT SUCCESS') : out('LOGOUT FAILURE');
});

bOpen.addEventListener('click', function () {
  if (ws) {
    return false;
  }
  // ws = new WebSocket('wss:\/\/signal-server.fly.dev\/echo');
  ws = new WebSocket('ws://localhost:8080/echo');
  ws.onopen = function () {
    out('OPEN');
  };
  ws.onclose = function () {
    out('CLOSE');
    ws = null;
  };
  ws.onmessage = function (evt) {
    out('RESPONSE: ' + evt.data);
  };
  ws.onerror = function (evt) {
    console.log(evt);
    out('ERROR');
  };
  return false;
});

bSend.addEventListener('click', function () {
  if (!ws) {
    return false;
  }
  out('SEND: ' + iInput.value);
  ws.send(iInput.value);
  return false;
});

bClose.addEventListener('click', function () {
  if (!ws) {
    return false;
  }
  ws.close();
  return false;
});
