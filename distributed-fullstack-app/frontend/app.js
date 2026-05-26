const BASE_URL = 'http://localhost:8080/api';
const consoleLog = document.getElementById('output');

function showResponse(title, data) {
    consoleLog.textContent = `=== ${title} ===\n${JSON.stringify(data, null, 2)}`;
}

const myClientID = crypto.randomUUID();
console.log("Meu ID de cliente único:", myClientID);

async function subscribeCategory() {
    const text = document.getElementById('category-payload').value;
    try {
        const res = await fetch(`${BASE_URL}/categories/subscribe`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ payload: text, client_id: myClientID }) // Envia a categoria e o ID
        });
        showResponse('Subscribe', await res.json());
    } catch (err) { showResponse('Error', err.message); }
}

async function unsubscribeCategory() {
    const text = document.getElementById('category-payload').value;
    try {
        const res = await fetch(`${BASE_URL}/categories/unsubscribe`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ payload: text, client_id: myClientID })
        });
        showResponse('Unsubscribe', await res.json());
    } catch (err) { showResponse('Error', err.message); }
}

async function registerPromotion() {
    const text = document.getElementById('prod-payload').value;
    try {
        const res = await fetch(`${BASE_URL}/promotions/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ payload: text })
        });
        showResponse('Register Promotion', await res.json());
    } catch (err) { showResponse('Error', err.message); }
}

async function voteInPromotion() {
    const name = document.getElementById('vote-name').value;
    const vote = document.getElementById('vote-type').value;

    const plainString = `${name} ${vote}`;

    try {
        const res = await fetch(`${BASE_URL}/promotions/vote`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ payload: plainString }) 
        });
        showResponse('Vote Status', await res.json());
    } catch (err) { showResponse('Error', err.message); }
}

async function listPromotions() {
    try {
        const res = await fetch(`${BASE_URL}/promotions/list`);
        showResponse('Promotions History', await res.json());
    } catch (err) { showResponse('Error', err.message); }
}

document.getElementById('btn-register').addEventListener('click', registerPromotion);
document.getElementById('btn-vote').addEventListener('click', voteInPromotion);
document.getElementById('btn-list').addEventListener('click', listPromotions);
document.getElementById('btn-subscribe').addEventListener('click', subscribeCategory);
document.getElementById('btn-unsubscribe').addEventListener('click', unsubscribeCategory);

