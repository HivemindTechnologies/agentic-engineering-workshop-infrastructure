/* HIVEMIND workshop theme behaviour — simplified fork of the Lambda World 2025
   deck's hivemind.js. Removed: the dual-presenter (IG/IB) initials indicator
   and its slidechanged handler. Kept: bracketed h1/h2 headings and the footer
   logo. Progress chrome: one bar over the whole deck, a ☕ per break. */

Reveal.on('ready', () => {
  // Bracket h1s (cyan) and h2s (red), except on the title slide.
  document.querySelectorAll('.reveal h1').forEach((h) => {
    if (h.dataset.bracketed || h.closest('.title-slide')) return;
    h.innerHTML =
      '<span style="color:#00c5db">[</span>' + h.innerHTML + '<span style="color:#00c5db">]</span>';
    h.dataset.bracketed = 'true';
  });
  document.querySelectorAll('.reveal h2').forEach((h) => {
    if (h.dataset.bracketed || h.closest('.title-slide')) return;
    h.innerHTML =
      '<span style="color:#f93030">[</span>' + h.innerHTML + '<span style="color:#f93030">]</span>';
    h.dataset.bracketed = 'true';
  });

  // Footer HIVEMIND logo (once).
  if (!document.querySelector('.reveal .footer-brand')) {
    const footer = document.createElement('div');
    footer.className = 'footer-brand';
    const img = document.createElement('img');
    img.src = 'assets/images/Logotype.svg';
    img.alt = 'HIVEMIND';
    footer.appendChild(img);
    document.querySelector('.reveal').appendChild(footer);
  }

  // Deck title in the footer (was top-right; M3a moves it off the header).
  if (!document.querySelector('.reveal .presentation-title')) {
    const label = document.createElement('div');
    label.className = 'presentation-title';
    label.textContent = document.title || 'HIVEMIND Workshop';
    document.querySelector('.reveal').appendChild(label);
  }

  initProgressChrome();
  updateProgressChrome(Reveal.getCurrentSlide());
  Reveal.on('slidechanged', (event) => {
    updateProgressChrome(event.currentSlide);
  });
});

function initProgressChrome() {
  if (document.querySelector('.reveal .deck-progress')) return;
  const root = document.createElement('div');
  root.className = 'deck-progress is-hidden';
  root.setAttribute('aria-hidden', 'true');
  root.innerHTML =
    '<span class="deck-progress-index"></span>' +
    '<div class="deck-progress-bar" role="presentation">' +
    // Full-width dashed track underneath; solid elapsed covers it as we advance.
    '<div class="deck-progress-track"></div>' +
    '<div class="deck-progress-elapsed"></div>' +
    '<div class="deck-progress-dots"></div>' +
    '</div>' +
    '<span class="deck-progress-title"></span>';
  document.querySelector('.reveal').appendChild(root);
}

// Global progress: the bar spans the half-day, with a ☕ where each break
// falls, so the room sees how far it is and when the next break comes.
// Reveal's own slide list gives the position; build-deck's
// data-section-title names the block under the bar. One dot per slide would
// overlap at ~80 slides, so only the current slide gets a dot.
function updateProgressChrome(slide) {
  const chrome = document.querySelector('.reveal .deck-progress');
  if (!chrome) return;

  const slides = Reveal.getSlides();
  const idx = slides.indexOf(slide);
  const count = slides.length;
  // Every slide but the title shows the bar: the opening counts toward session 1.
  const isTitle = slide && (slide.classList.contains('title-slide') || slide.querySelector('.title-slide'));
  if (!slide || isTitle || idx < 0 || count < 2) {
    chrome.classList.add('is-hidden');
    chrome.setAttribute('aria-hidden', 'true');
    return;
  }
  chrome.classList.remove('is-hidden');
  chrome.setAttribute('aria-hidden', 'false');

  // Every session lasts 50 minutes however many slides it holds, so each gets
  // an equal share of the bar: the breaks sit at 1/4, 2/4, 3/4, and slides
  // spread evenly inside their session.
  const anchors = [0];
  slides.forEach((s, i) => { if (s.classList.contains('break-slide')) anchors.push(i); });
  anchors.push(count - 1);
  const at = (i) => {
    let k = 0;
    while (k < anchors.length - 2 && i >= anchors[k + 1]) k++;
    const span = anchors[k + 1] - anchors[k] || 1;
    return ((k + (i - anchors[k]) / span) / (anchors.length - 1)) * 100 + '%';
  };
  chrome.querySelector('.deck-progress-index').textContent = idx + 1 + ' / ' + count;
  // The block this slide belongs to ("Basics"), so a room that lost track
  // sees where it is without going back to the divider slide.
  chrome.querySelector('.deck-progress-title').textContent = slide.dataset.sectionTitle || '';
  // Dashed track is always full-width in CSS; only the solid cover grows.
  chrome.querySelector('.deck-progress-elapsed').style.width = at(idx);

  // The deck does not change while it is shown: build the markers once.
  const marks = chrome.querySelector('.deck-progress-dots');
  if (!marks.dataset.built) {
    marks.dataset.built = 'true';
    slides.forEach((s, i) => {
      if (!s.classList.contains('break-slide')) return;
      const cup = document.createElement('span');
      cup.className = 'deck-progress-break';
      cup.textContent = '☕';
      cup.style.left = at(i);
      cup.dataset.slide = String(i);
      marks.appendChild(cup);
    });
    const dot = document.createElement('span');
    dot.className = 'deck-progress-dot is-current';
    marks.appendChild(dot);
  }
  marks.querySelector('.deck-progress-dot').style.left = at(idx);
  marks.querySelectorAll('.deck-progress-break').forEach((cup) => {
    cup.classList.toggle('is-elapsed', Number(cup.dataset.slide) < idx);
  });
}
