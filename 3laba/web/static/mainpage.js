document.getElementById('download').addEventListener('click', download);

async function downloadUpdated() {
    const rows = document.querySelectorAll('#dataTable tbody tr');
    const data = Array.from(rows).map(row => {
        const cells = row.querySelectorAll('td');
        const rawDate = cells[2].innerText.trim().replace(/\./g, '-');
        return {
            Name: cells[0].innerText.trim(),
            NameOfWork: cells[1].innerText.trim(),
            Date: rawDate + "T00:00:00Z",
            Type: cells[3].innerText.trim()
        };
    });

    const response = await fetch('/save', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });

    if (!response.ok) {
        alert(response.status + '' + response.statusText);
        return;
    }
    const text = await response.text();
    const textarea = document.getElementById('resultArea');
    textarea.value = text;
}

function newString(){
    const dataTable = document.getElementById('dataTable').getElementsByTagName('tbody')[0];
    const newRows = dataTable.insertRow();
    for (let i = 0; i < 4; i++) {
        const cell = newRows.insertCell();
        cell.contentEditable = true;
        cell.textContent = '';
    }
}

async function download() {

    const outputArea = document.getElementById('resultArea');
    const data = outputArea.value;

    if (!outputArea.value) {
        alert('Нет данных для скачивания');
        return;
    }

    const blob = new Blob([data], {type: 'text/plain;charset=UTF-8'});
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'output.txt';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}
function clearTable() {
    const tbody = document.querySelector('#dataTable tbody');
    tbody.innerHTML = '';
}
