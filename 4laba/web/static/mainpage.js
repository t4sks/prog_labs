const downloadBtn = document.getElementById('download');
if (downloadBtn) {
    downloadBtn.addEventListener('click', download);
}

function collectTableData() {
    const rows = document.querySelectorAll('#dataTable tbody tr');

    return Array.from(rows).map(row => {
        const cells = row.querySelectorAll('td');
        const rawDate = cells[2].innerText.trim().replace(/\./g, '-');

        return {
            Name: cells[0].innerText.trim(),
            NameOfWork: cells[1].innerText.trim(),
            Date: rawDate ? rawDate + 'T00:00:00Z' : '',
            Type: cells[3].innerText.trim()
        };
    });
}

function renderTable(data) {
    const tbody = document.querySelector('#dataTable tbody');
    tbody.innerHTML = '';

    data.forEach(item => {
        const row = document.createElement('tr');

        let formattedDate = '';
        if (item.Date) {
            const parsedDate = new Date(item.Date);
            if (!isNaN(parsedDate.getTime())) {
                formattedDate = parsedDate.toISOString().slice(0, 10).replace(/-/g, '.');
            }
        }

        row.innerHTML = `
            <td contenteditable="true">${item.Name ?? ''}</td>
            <td contenteditable="true">${item.NameOfWork ?? ''}</td>
            <td contenteditable="true">${formattedDate}</td>
            <td contenteditable="true">${item.Type ?? ''}</td>
        `;

        tbody.appendChild(row);
    });
}

function newString() {
    const dataTable = document.getElementById('dataTable').getElementsByTagName('tbody')[0];
    const newRow = dataTable.insertRow();

    for (let i = 0; i < 4; i++) {
        const cell = newRow.insertCell();
        cell.contentEditable = true;
        cell.textContent = '';
    }
}

function clearTable() {
    const tbody = document.querySelector('#dataTable tbody');
    tbody.innerHTML = '';
}

async function downloadUpdated() {
    const data = collectTableData();

    const response = await fetch('/save', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });

    if (!response.ok) {
        alert(response.status + ' ' + response.statusText);
        return;
    }

    const text = await response.text();
    const textarea = document.getElementById('resultArea');
    textarea.value = text;
}

async function download() {
    const outputArea = document.getElementById('resultArea');
    const data = outputArea.value;

    if (!data) {
        alert('Нет данных для скачивания');
        return;
    }

    const blob = new Blob([data], { type: 'text/plain;charset=UTF-8' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');

    a.href = url;
    a.download = 'output.txt';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);

    URL.revokeObjectURL(url);
}

async function applyCommands() {
    const fileInput = document.getElementById('commandsFile');
    const errorArea = document.getElementById('ErrorArea');
    const resultArea = document.getElementById('resultArea');

    if (!fileInput.files.length) {
        alert('Выберите файл команд');
        return;
    }

    const works = collectTableData();

    const formData = new FormData();
    formData.append('commands', fileInput.files[0]);
    formData.append('data', JSON.stringify(works));

    const response = await fetch('/apply', {
        method: 'POST',
        body: formData
    });

    if (!response.ok) {
        alert(await response.text());
        return;
    }

    const updatedWorks = await response.json();

    const validWorks = [];
    const errors = [];

    updatedWorks.forEach(item => {
        const nameEmpty = !item.Name || item.Name.trim() === '';

        if (nameEmpty) {
            if (item.Error && item.Error.trim() !== '') {
                errors.push(item.Error);
            } else {
                errors.push('Некорректный объект без текста ошибки');
            }
            return;
        }

        validWorks.push(item);
    });

    renderTable(validWorks);
    errorArea.value = errors.join('\n');
    resultArea.value = '';
}