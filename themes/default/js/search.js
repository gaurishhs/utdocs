(async function () {
    async function fetchSearchIndex() {
        const response = await fetch('/search/index.json');
        return response.json();
    }
    var searchIndex = await fetchSearchIndex()

    const fuse = new Fuse(searchIndex, {
        keys: ['Title','Content'],
        minMatchCharLength: 3,
        includeMatches: true,
        shouldSort: true,
        threshold: 0.1,
        location: 0,
        distance: 100,
    });
    console.log('Fuse.js Initialized');

    // https://stackoverflow.com/questions/67352406/fuse-js-includematches-highlighting-how-does-it-work
    const highlight = (fuseSearchResult, highlightClassName = 'highlight') => {
        const set = (obj, path, value) => {
            const pathValue = path.split('.');
            let i;
      
            for (i = 0; i < pathValue.length - 1; i++) {
              obj = obj[pathValue[i]];
            }
      
            obj[pathValue[i]] = value;
        };
      
        const generateHighlightedText = (inputText, regions = []) => {
          let content = '';
          let nextUnhighlightedRegionStartingIndex = 0;
      
          regions.forEach(region => {
            const lastRegionNextIndex = region[1] + 1;
      
            content += [
              inputText.substring(nextUnhighlightedRegionStartingIndex, region[0]),
              `<span class="${highlightClassName}">`,
              inputText.substring(region[0], lastRegionNextIndex),
              '</span>',
            ].join('');
      
            nextUnhighlightedRegionStartingIndex = lastRegionNextIndex;
          });
      
          content += inputText.substring(nextUnhighlightedRegionStartingIndex);
      
          return content;
        };
      
        return fuseSearchResult
          .filter(({ matches }) => matches && matches.length)
          .map(({ item, matches }) => {
            const highlightedItem = { ...item };
      
            matches.forEach((match) => {
              set(highlightedItem, match.key, generateHighlightedText(match.value, match.indices));
            });
      
            return highlightedItem;
          });
    };

    var searchInput = document.getElementById('search-input');

    searchInput.addEventListener('keyup', function (e) {
        var searchResults = document.getElementById('search-results');
        var searchValue = e.target.value;
        var results = highlight(fuse.search(searchValue));
        var html = '';
        if (results) {
            results.map((result) => {
                html += `
                    <li>
                        <a class="search-result" href="${result.Url}">
                            <span>${result.Title}</span>
                            <p>${result.Content}</p>
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
        searchInput.focus();
    }

    var searchModalCloseBtn = document.getElementById('search-modal-close');
    searchModalCloseBtn.onclick = function () {
        searchContainer.style.display = 'none';
    }
})();