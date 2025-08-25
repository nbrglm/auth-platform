function setTheme() {
  let theme = window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light"
  document.documentElement.setAttribute("data-theme", theme);
  Alpine.store('darkMode', theme === "dark")
}

function setupImageEditor() {
  Alpine.store('imageEditor', {isLoaded: false, setLoaded(val) { this.isLoaded = val }});
}

function setupAlerts() {
  Alpine.store('alerts', {
    kindError: "error",
    kindSuccess: "success",
    kindWarn: "warn",
    kindInfo: "info",

    items: [],

    add(message, kind = "info", duration = 3000) {
      const id = Date.now()
      this.items.push({id, message, kind})
      if (duration) setTimeout(() => {this.remove(id)}, duration)
      return id;
    },

    remove(id) {
      this.items = this.items.filter(i => i.id !== id)
    }
  })
}

function setupActions() {
  Alpine.store('actions', {
    show: false,
    message: "",
    title: "",
    confirmButtonText: "",
    cancelButtonText: "",
    resolve: null,
    reject: null,
    id: "actionsModalDialog",

    ask(title, message, confirmButtonText = "Confirm", cancelButtonText = "Cancel") {
      this.title = title;
      this.message = message;
      this.confirmButtonText = confirmButtonText;
      this.cancelButtonText = cancelButtonText;

      return new Promise((resolve, reject) => {
        this.resolve = resolve;
        this.reject = reject;
      })
    },

    confirm() {
      this.resolve?.(true);
      this.reset();
    },

    cancel() {
      // We need to null check reject, since cancel is also called
      // we then dialog receives a close event,
      // and this close event is fired when the dialog closes,
      // so it might be that the dialog was closed with the confirm method above,
      // it will still call the cancel method, hence, we null check reject to make sure
      // the dialog was not closed with confirm first.
      this.reject?.(false);
      this.reset();
    },

    reset() {
      this.show = false;
      this.reject = null;
      this.resolve = null;
      this.title = "";
      this.message = "";
    }
  })
}

function initializeState() {
  setTheme();
  Alpine.store('sidebarState', {isOpen: false, setOpen(val) {
    this.isOpen = val;
  }})

  // setup the image editor store
  setupImageEditor();

  // Setup alerts (toasts)
  setupAlerts();

  // Setup actions (confirmation dialogs)
  setupActions();
}

window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", setTheme);

document.addEventListener("alpine:init", () => {
  initializeState();
})