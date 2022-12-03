(async function () {
    async function fetchSearchIndex() {
        const response = await fetch('/search/index.json');
        return response.json();
    }
    var searchIndex = await fetchSearchIndex()

    const fuse = new Fuse(searchIndex, {
        keys: ['Title', 'Url'],
        minMatchCharLength: 3,
    });
    console.log('Fuse.js Initialized');

    var searchInput = document.getElementById('search-input');

    searchInput.addEventListener('keyup', function (e) {
        var searchResults = document.getElementById('search-results');
        var searchValue = e.target.value;
        var results = fuse.search(searchValue);
        var html = '';
        if (results) {
            results.map((result) => {
                html += `
                    <li>
                        <a class="search-result" href="${result.item.Url}">
                            ${result.item.Title}
                        </a>
                    </li>
                `;
            });
        }
        if (!results.length) {
            html = '<center><span style="color: #7f8497; font-size: 0.9em; text-align: center; padding-top: 10px; padding-bottom: 10px;">No results found. Make sure to have atleast 3 characters.</span></center>';
        }
        searchResults.innerHTML = html;
    });
    var searchContainer = document.getElementById('search-container');
    var searchBtn = document.getElementById('search_btn');
    searchBtn.onclick = function () {
        searchContainer.style.display = 'block';
    }

    var searchModalCloseBtn = document.getElementById('search-modal-close');
    searchModalCloseBtn.onclick = function () {
        searchContainer.style.display = 'none';
    }
})();