var versions = [
  {path: "master", text: "Experimental", selected: false },
  {path: "v1.10", text: "v1.10 (RC)", selected: false },
  {path: "", text: "v1.9 Latest", selected: false },
  {path: "v1.8", text: "v1.8", selected: true },
];




function createBanner(parentElem, versions) {
  if (!parentElem || window.location.host !== "doc.traefik.io") {
    return;
  }

  const products = {
    traefik: {
      color: '#2aa2c1',
      backgroundColor: '#2aa2c11a',
      fullName: 'Traefik Proxy',
    },
    'traefik-enterprise': {
      color: '#337fe6',
      backgroundColor: '#337fe61a',
      fullName: 'Traefik Enterprise',
    },
    'traefik-mesh': {
      color: '#be46dd',
      backgroundColor: '#be46dd1a',
      fullName: 'Traefik Mesh',
    },
  }

  const [,productName] = window.location.pathname.split('/');
  const currentProduct = products[productName];
  const currentVersion = versions.find(v => v.selected);
  const preExistentBanner = document.getElementById('outdated-doc-banner');

  if (!currentProduct || !currentVersion || !!preExistentBanner) return;

  const cssCode = `
    #obsolete-banner {
      display: flex;
      width: 100%;
      align-items: center;
      justify-content: center;
      max-width: 1274px;
      margin: 0 auto;
    }
    #obsolete-banner .obsolete-banner-content {
      display: flex;
      align-items: center;
      height: 40px;
      margin: 24px;
      padding: 11px 16px;
      border-radius: 8px;
      background-color: ${currentProduct.backgroundColor};
      gap: 16px;
      font-family: Rubik, sans-serif;
      font-size: 14px;
      color: ${currentProduct.color};
      box-sizing: border-box;
      width: 100%;
    }
    #obsolete-banner .obsolete-banner-content strong { font-weight: bold; }
    #obsolete-banner .obsolete-banner-content a { color: ${currentProduct.color}; text-decoration: none; }
    #obsolete-banner .obsolete-banner-content a:hover { text-decoration: underline; }
    #obsolete-banner .obsolete-banner-content p { margin: 0; }
  `

  const banner = document.createElement('div');
  banner.id = 'obsolete-banner';
  banner.innerHTML = `
    <div class="obsolete-banner-content">
      <strong>OLD VERSION</strong>
      <p>
        You're looking at documentation for ${currentProduct.fullName} ${currentVersion.text}.
        <a href="/${productName}">Click here to see the latest version. →</a>
      </p>
    </div>
  `;

  // Append HTML
  parentElem.prepend(banner);

  // Append Styling
  const [head] = document.getElementsByTagName("head");
  if (!document.getElementById("obsolete-banner-style")) {
    const styleElem = document.createElement("style");
    styleElem.id = "obsolete-banner-style";

    if (styleElem.styleSheet) {
      styleElem.styleSheet.cssText = cssCode;
    } else {
      styleElem.appendChild(document.createTextNode(cssCode));
    }

    head.appendChild(styleElem);
  }
}

function addBannerMaterial(versions) {
  const elem = document.querySelector('body > div.md-container');
  createBanner(elem, versions)
}

function addBannerUnited() {
  const elem = document.querySelector('body > div.container');
  createBanner(elem, versions)
}

// Material theme

function addMaterialMenu(elt, versions) {
  const current = versions.find(function (value) {
    return value.selected;
  })

  const rootLi = document.createElement('li');
  rootLi.classList.add('md-nav__item');
  rootLi.classList.add('md-nav__item--nested');

  const input = document.createElement('input');
  input.classList.add('md-toggle');
  input.classList.add('md-nav__toggle');
  input.setAttribute('data-md-toggle', 'nav-10000000');
  input.id = "nav-10000000";
  input.type = 'checkbox';

  rootLi.appendChild(input);

  const lbl01 = document.createElement('label');
  lbl01.classList.add('md-nav__link');
  lbl01.setAttribute('for', 'nav-10000000');

  const spanTitle01 = document.createElement('span');
  spanTitle01.classList.add('md-nav__item-title');
  spanTitle01.textContent = current.text+ " ";

  lbl01.appendChild(spanTitle01);

  const spanIcon01 = document.createElement('span');
  spanIcon01.classList.add('md-nav__icon');
  spanIcon01.classList.add('md-icon');

  lbl01.appendChild(spanIcon01);

  rootLi.appendChild(lbl01);

  const nav = document.createElement('nav')
  nav.classList.add('md-nav');
  nav.setAttribute('data-md-component','collapsible');
  nav.setAttribute('aria-label', current.text);
  nav.setAttribute('data-md-level','1');

  rootLi.appendChild(nav);

  const lbl02 = document.createElement('label');
  lbl02.classList.add('md-nav__title');
  lbl02.setAttribute('for', 'nav-10000000');
  lbl02.textContent = current.text + " ";

  const spanIcon02 = document.createElement('span');
  spanIcon02.classList.add('md-nav__icon');
  spanIcon02.classList.add('md-icon');

  lbl02.appendChild(spanIcon02);

  nav.appendChild(lbl02);

  const ul = document.createElement('ul');
  ul.classList.add('md-nav__list');
  ul.setAttribute('data-md-scrollfix','');

  nav.appendChild(ul);

  for (let i = 0; i < versions.length; i++) {
    const li = document.createElement('li');
    li.classList.add('md-nav__item');

    ul.appendChild(li);

    const a = document.createElement('a');
    a.classList.add('md-nav__link');
    if (versions[i].selected) {
      a.classList.add('md-nav__link--active');
    }
    a.href = window.location.protocol + "//" + window.location.host + "/";
    if (window.location.host === "doc.traefik.io") {
      a.href = a.href + window.location.pathname.split('/')[1] + "/";
    }
    if (versions[i].path) {
      a.href = a.href + versions[i].path + "/";
    }
    a.title = versions[i].text;
    a.text = versions[i].text;

    li.appendChild(a);
  }

  elt.appendChild(rootLi);
}

// United theme

function addMenu(elt, versions){
  const li = document.createElement('li');
  li.classList.add('md-nav__item');
  li.style.cssText = 'padding-top: 1em;';

  const select = document.createElement('select');
  select.classList.add('md-nav__link');
  select.style.cssText = 'background: white;border: none;color: #00BCD4;-webkit-border-radius: 5px;-moz-border-radius: 5px;border-radius: 5px;overflow: hidden;padding: 0.1em;'
  select.setAttribute('onchange', 'location = this.options[this.selectedIndex].value;');

  for (let i = 0; i < versions.length; i++) {
    let opt = document.createElement('option');
    opt.value = window.location.protocol + "//" + window.location.host + "/";
    if (window.location.host === "doc.traefik.io") {
      opt.value = opt.value + window.location.pathname.split('/')[1] + "/"
    }
    if (versions[i].path) {
      opt.value = opt.value + versions[i].path + "/"
    }
    opt.text = versions[i].text;
    opt.selected = versions[i].selected;
    select.appendChild(opt);
  }

  li.appendChild(select);
  elt.appendChild(li);
}


const unitedSelector = 'div.navbar.navbar-default.navbar-fixed-top div.container div.navbar-collapse.collapse ul.nav.navbar-nav.navbar-right';
const materialSelector = 'div.md-container main.md-main div.md-main__inner.md-grid div.md-sidebar.md-sidebar--primary div.md-sidebar__scrollwrap div.md-sidebar__inner nav.md-nav.md-nav--primary ul.md-nav__list';

let elt = document.querySelector(materialSelector);
if (elt) {
  addMaterialMenu(elt, versions);
  addBannerMaterial(versions);
} else {
  const elt = document.querySelector(unitedSelector);
  addMenu(elt, versions);
  addBannerUnited(versions);
}
