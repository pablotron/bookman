(() => {
  "use strict";

  // get element by ID, get elements by selector
  const get = (id) => document.getElementById(id);
  const qsa = (s) => document.querySelectorAll(s);
  const on = (el, ev, fn) => el.addEventListener(ev, fn);

  // cache search field and rows element
  const field = get('q'),
        books = get('books'),
        upload = get('upload');

  // html escape
  const h = (v) => {
    return String(v).replaceAll('&', '&amp;')
      .replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;')
      .replaceAll("'", '&apos;')
      .replaceAll('"', '&quot;');
  };

  // templates
  const T = {
    // result template
    item: (row) => `
      <a
        href='/book/${h(row.id)}'
        class='panel-block'
        title='${h(row.name)}, by ${h(row.author)}'
        aria-label='${h(row.name)}, by ${h(row.author)}'
        data-id='${h(row.id)}'
        data-name='${h(row.name)}'
        data-author='${h(row.author)}'
        data-rank='${h(row.rank)}'
      >
        <span class='edit-book'>
          <svg xmlns='http://www.w3.org/2000/svg' width='16' height='16' fill='currentColor' class='bi bi-pencil-square' viewBox='0 0 16 16'>
            <path d='M15.502 1.94a.5.5 0 0 1 0 .706L14.459 3.69l-2-2L13.502.646a.5.5 0 0 1 .707 0l1.293 1.293zm-1.75 2.456-2-2L4.939 9.21a.5.5 0 0 0-.121.196l-.805 2.414a.25.25 0 0 0 .316.316l2.414-.805a.5.5 0 0 0 .196-.12l6.813-6.814z'/>
            <path fill-rule='evenodd' d='M1 13.5A1.5 1.5 0 0 0 2.5 15h11a1.5 1.5 0 0 0 1.5-1.5v-6a.5.5 0 0 0-1 0v6a.5.5 0 0 1-.5.5h-11a.5.5 0 0 1-.5-.5v-11a.5.5 0 0 1 .5-.5H9a.5.5 0 0 0 0-1H2.5A1.5 1.5 0 0 0 1 2.5v11z'/>
          </svg>
        </span>

        ${h(row.name)}, by ${h(row.author)}
      </a>
    `,

    // no match template
    none: () => `
      <div class='panel-block'>
        No matching results.
      </div>
    `,

    // list template
    list: (rows) => rows.map((row) => T.item(row)).join(''),
  };

  const refresh = () => {
    const new_q = field.value;
    const old_q = books.dataset.q || '';

    if (new_q === '' || new_q !== old_q) {
      // build url
      const url = 'api/search?' + (new URLSearchParams({ q: new_q })).toString();

      fetch(url).then((r) => r.json()).then((r) => {
        // cache query string
        books.dataset.q = q;

        // refresh list
        books.innerHTML = (r.length > 0) ? T.list(r): T.none();
      });
    }
  };

  on(document, 'DOMContentLoaded', () => {
    let t = null;

    // search field handler
    on(field, 'keydown', () => {
      // clear old timeout
      if (t !== null) {
        clearTimeout(t);
        t = null;
      }

      // refresh list after 200ms
      t = setTimeout(refresh, 200);
    });

    // edit btn handler
    on(get('books'), 'click', (ev) => {
      if (ev.target.closest('.edit-book')) {
        // get book data
        const data = ev.target.closest('a').dataset;

        // populate edit dialog
        get('edit-save-btn').dataset.id = data.id;
        get('edit-name').value = data.name;
        get('edit-author').value = data.author;

        // show edit dialog
        get('edit-dialog').classList.add('is-active');

        // stop event
        ev.preventDefault();
        return false;
      }
    });

    on(get('edit-save-btn'), 'click', (ev) => {
      // build form data
      const data = new FormData();
      data.append('id', get('edit-save-btn').dataset.id);
      data.append('name', get('edit-name').value);
      data.append('author', get('edit-author').value);

      // send request
      fetch('api/edit', {
        method: 'POST',
        body: data,
      }).then((r) => {
        if (!r.ok) {
          alert('edit failed');

          /* TODO: r.text().then(err => {
            console.log(err);
            alert(err);
          }); */

          return;
        }

        // hide dialog, refresh list
        get('edit-dialog').classList.remove('is-active');
        refresh();
      });

      // stop event
      ev.preventDefault();
      ev.stopPropagation();
      return false;
    });

    // upload btn handler
    on(get('upload-btn'), 'click', () => {
      // show upload dialog
      upload.click();
    });

    // upload dialog handler
    on(upload, 'change', () => {
      const files = upload.files;
      if (files.length == 0) {
        return;
      }
      console.log(files);

      // build post body
      let data = new FormData();
      for (let f of files) {
        data.append('file', f);
      }

      // fetch files
      fetch('api/upload', {
        method: 'POST',
        body: data,
      }).then((r) => {
        if (r.ok) {
          refresh();
        } else {
          alert('upload failed');
        }
      });
    });

    // modal handlers
    // (adapted from https://bulma.io/documentation/components/modal/)
    const modal_show = (e) => e.classList.add('is-active'),
          modal_hide = (e) => e.classList.remove('is-active'),
          hide_all_modals = () => (qsa('.modal') || []).forEach((e) => modal_hide(e));

		// Add a click event handlers to close parent modal
		(qsa('.modal-background, .modal-close, .modal-card-head .delete, .modal-card-foot') || []).forEach((e) => {
      const modal = e.closest('.modal');
			on(e, 'click', () => modal_hide(modal));
		});

		// Add a keyboard event to close all modals
		on(document, 'keydown', (ev) => {
      // check for escape key
			if ((ev || window.event).keyCode === 27) {
				hide_all_modals();
			}
		});
  });

  // load initial list
  refresh();
})();
