class Alert extends HTMLElement {
    constructor() { super(); }
    render() {
        this.innerHTML = `
            <div class="alert">
                <div class="d-flex">
                    <div class="flex-shrink-0">
                        ${this.resolveIcon(this.getAttribute('type'))}
                    </div>
                    <div style="margin-left: .75rem;">
                        <h3 style="font-size: 1rem; color: ${this.resolveColor(this.getAttribute('type'))}; font-weight: 500;" class="text-base font-medium">${this.resolveText(this.getAttribute('type'))}</h3>
                        <div class="mt-2 text-sm-start">
                            ${this.getAttribute('message')}
                        </div>
                    </div>
                </div>
            </div>
        `;
    }
    resolveText(type) {
        switch (type) {
            case 'error':
                return 'Error';
            case 'success':
                return 'Success';
            case 'info':
                return 'Info';
            case 'warn':
                return 'Warning';
        }
    }

    resolveColor(type) {
        switch (type) {
            case 'warn':
                return 'rgb(253 173 65 / 1)';
            case 'error':
                return 'rgb(254 107 136 / 1);'
            case 'success':
                return 'rgb(0 224 204 / 1);'
            case 'info':
                return 'rgb(85 118 249 / 1);'
        }
    }

    resolveIcon(type) {
        switch (type) {
            case 'warn':
                return `<svg xmlns="http://www.w3.org/2000/svg" class="
                icon icon-tabler
                text-warning
                icon-tabler-alert-triangle" style="color: ${this.resolveColor(type)}" width="20" height="20" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" fill="none"
              stroke-linecap="round" stroke-linejoin="round">
              <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
              <path d="M12 9v2m0 4v.01"></path>
              <path d="M5 19h14a2 2 0 0 0 1.84 -2.75l-7.1 -12.25a2 2 0 0 0 -3.5 0l-7.1 12.25a2 2 0 0 0 1.75 2.75">
              </path>
            </svg>`
                break;
            case 'info':
                return `<svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-info-circle text-oyster"
                width="20" style="color: ${this.resolveColor(type)}" height="20" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none"
                stroke-linecap="round" stroke-linejoin="round">
                <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                <circle cx="12" cy="12" r="9"></circle>
                <line x1="12" y1="8" x2="12.01" y2="8"></line>
                <polyline points="11 12 12 12 12 16 13 16"></polyline>
              </svg>`
              break;
            case 'success':
                return `<svg xmlns="http://www.w3.org/2000/svg" class="
                icon icon-tabler
                icon-tabler-circle-check
                " width="20" style="color: ${this.resolveColor(type)}" height="20" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" fill="none"
                stroke-linecap="round" stroke-linejoin="round">
                <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                 <circle cx="12" cy="12" r="9"></circle>
              <path d="M9 12l2 2l4 -4"></path>
             </svg>`
                break;
            case 'error':
                return `<svg xmlns="http://www.w3.org/2000/svg" class="
                icon icon-tabler
                text-psycho
                icon-tabler-alert-octagon
                text-brink
              " width="20" height="20" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" fill="none"
              stroke-linecap="round" style="color: ${this.resolveColor(type)}" stroke-linejoin="round">
              <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
              <path
                d="M8.7 3h6.6c.3 0 .5 .1 .7 .3l4.7 4.7c.2 .2 .3 .4 .3 .7v6.6c0 .3 -.1 .5 -.3 .7l-4.7 4.7c-.2 .2 -.4 .3 -.7 .3h-6.6c-.3 0 -.5 -.1 -.7 -.3l-4.7 -4.7c-.2 -.2 -.3 -.4 -.3 -.7v-6.6c0 -.3 .1 -.5 .3 -.7l4.7 -4.7c.2 -.2 .4 -.3 .7 -.3z">
              </path>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>`
        }
    }

    connectedCallback() {
        if (!this.rendered) {
            this.render();
            this.rendered = true;
        }
    }

    static get observedAttributes() {
        return ['type'];
    }

    attributeChangedCallback(name, oldValue, newValue) {
        this.render();
    }

}


customElements.define('u-alert', Alert);