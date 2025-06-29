class SearchModal {
    constructor() {
        this.messageIndex = null;
        this.focusedIndex = 0;
        this.init();
    }

    init() {
        document.getElementById('search').addEventListener('keydown', (event) => {
            const resultsList = document.getElementById('search-results').querySelectorAll('li');

            if (event.key === 'Enter') {
                event.preventDefault();
                this.handleEnterKey();
            }

            if (event.key === 'ArrowDown') {
                event.preventDefault();
                this.focusedIndex = (this.focusedIndex + 1) % resultsList.length;
                this.updateFocusedItem();
            } else if (event.key === 'ArrowUp') {
                event.preventDefault();
                this.focusedIndex = (this.focusedIndex - 1 + resultsList.length) % resultsList.length;
                this.updateFocusedItem();
            }

            if (event.key === 'Escape') {
                this.close();
                closeMove();
            }
        });
    }

    open(text = '', messageIndex = null) {
        this.messageIndex = messageIndex;

        document.getElementById('search').style.display = 'block';
        const inputField = document.getElementById('search-input');
        inputField.value = text;
        inputField.focus();

        this.focusedIndex = 0;
        const goToFileResults = document.getElementById('search-results');
        goToFileResults.innerHTML = '';

        if (text === '') {
            this.loadRecentFiles();
        } else {
            search();
        }
    }

    close() {
        document.getElementById('search').style.display = 'none';
        this.messageIndex = null;
    }

    showResults(results) {
        const list = document.getElementById('search-results');
        list.innerHTML = ''; // Clear previous results

        results.forEach(({dir, filename}, index) => {
            if (filename === CONFIG_FILENAME) {
                return;
            }

            const listItem = document.createElement('li');
            let title = filename.replace(/\.md$/, '')
            if (dir !== '') {
                listItem.textContent = `${dir}/${title}`;
            } else {
                listItem.textContent = title;
            }
            listItem.setAttribute('data-path', `${dir}/${filename}`);
            listItem.setAttribute('data-index', index);

            listItem.onclick = () => this.handleItemClick(dir, filename);

            listItem.onmouseenter = () => {
                document.querySelectorAll('#search-results li').forEach(li => li.classList.remove('focused'));
                listItem.classList.add('focused');
                this.focusedIndex = index;
            };
            list.appendChild(listItem);
        });

        this.focusedIndex = 0;
        this.updateFocusedItem();
    }

    handleItemClick(dir, filename) {
        console.log('handling');
        if (this.messageIndex !== null) {
            sendCmd('mf', [filename, this.messageIndex.toString()]);
            this.close();
        } else {
            // Default behavior
            openEditor(!isChat);
            openFile(dir, filename);
            this.close();
        }
    }

    handleEnterKey() {
        const resultsList = document.getElementById('search-results').querySelectorAll('li');
        if (resultsList[this.focusedIndex]) {
            const [dir, filename] = resultsList[this.focusedIndex].getAttribute('data-path').split('/');
            this.handleItemClick(dir, filename);
        }
    }

    updateFocusedItem() {
        const resultsList = document.getElementById('search-results').querySelectorAll('li');
        document.querySelectorAll('#search-results li').forEach(li => li.classList.remove('focused'));
        resultsList.forEach((item, index) => {
            if (index === this.focusedIndex) {
                item.classList.add('focused');
                item.scrollIntoView({block: 'nearest'});
            } else {
                item.classList.remove('focused');
            }
        });
    }

    loadRecentFiles() {
        let results = [];
        for (const dir of Object.keys(excludeDirs(SYSTEM_DIRS))) {
            for (const filename of Object.keys(files[dir])) {
                results.push({
                    dir, filename, lastModified: files[dir][filename].lastModified,
                });
            }
        }

        results = results
            .sort((a, b) => b.lastModified - a.lastModified)
            .slice(0, 8);

        this.showResults(results);
    }
}

const searchModal = new SearchModal();