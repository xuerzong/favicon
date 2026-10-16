import { debounce } from "./debounce";
import "./style.css";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/";
const DEFAULT_DOMAIN = "example.com";

const input = document.getElementById("url-input") as HTMLInputElement;
const wrapper = document.getElementById("input-wrapper") as HTMLLabelElement;
const prefix = document.querySelector(".input-prefix") as HTMLSpanElement;
const inputSuffix = document.querySelector(".input-suffix") as HTMLSpanElement;
const preview = document.getElementById("favicon-preview") as HTMLImageElement;
const loading = document.getElementById("loading") as HTMLDivElement;

prefix.textContent = API_BASE_URL;

const handleUrlChange = debounce((nextUrl: string) => {
  preview.src = `${API_BASE_URL}${nextUrl || DEFAULT_DOMAIN}`;
  loading.hidden = true;
}, 1000);

input.addEventListener("focus", () => {
  wrapper.dataset.focused = "true";
});

input.addEventListener("blur", () => {
  delete wrapper.dataset.focused;
});

input.addEventListener("input", () => {
  loading.hidden = false;
  handleUrlChange(input.value);
});

inputSuffix.addEventListener("click", (e) => {
  e.preventDefault();
  window.open(`${API_BASE_URL}${input.value || "example.com"}`, "_blank");
});

type FormItem = {
  id: string;
  label: string;
  name: string;
  type: "numberInput" | "select";
  defaultValue: string;
  options?: { text: string; value: string }[];
};

const formItems: FormItem[] = [
  {
    id: "size-input",
    label: "Size",
    name: "size",
    defaultValue: "32",
    type: "numberInput",
  },
  {
    id: "format-select",
    label: "Format",
    name: "format",
    type: "select",
    defaultValue: "png",
    options: ["png", "webp"].map((op) => ({
      text: op,
      value: op,
    })),
  },
  {
    id: "rotate-select",
    label: "Rotate",
    name: "rotate",
    type: "select",
    defaultValue: "0",
    options: [0, 90, 180, 270].map(String).map((op) => ({
      text: op,
      value: op,
    })),
  },
];

function renderQueriesForm() {
  const queriesEle = document.getElementById("queries")!;
  const queriesFormEle = document.createElement("form");
  queriesFormEle.classList.add("form");

  queriesEle.appendChild(queriesFormEle);

  formItems.forEach((formItem) => {
    const container = document.createElement("div");
    container.classList.add("form-item");

    const label = document.createElement("label");
    label.classList.add("form-label");
    label.setAttribute("htmlFor", formItem.name);
    const labelSpan = document.createElement("span");
    labelSpan.textContent = formItem.label;
    label.appendChild(labelSpan);

    let ele: HTMLElement | null = null;
    switch (formItem.type) {
      case "select": {
        ele = document.createElement("select");
        const selectEle: HTMLSelectElement = ele as HTMLSelectElement;
        selectEle.setAttribute("id", formItem.id);
        formItem.options?.forEach((op) => {
          selectEle.add(new Option(op.text, op.value));
        });
        selectEle.value = formItem.defaultValue;
        selectEle.addEventListener("change", (e) => {
          updateQueryString({ [formItem.name]: (e.target as any).value });
        });
        break;
      }
      case "numberInput": {
        ele = document.createElement("input");
        const inputEle = ele as HTMLInputElement;
        inputEle.setAttribute("id", formItem.id);
        inputEle.setAttribute("type", "number");
        inputEle.value = formItem.defaultValue;
        inputEle.addEventListener("change", (e) => {
          updateQueryString({ [formItem.name]: (e.target as any).value });
        });
        break;
      }
      default:
        break;
    }

    container.appendChild(label);
    if (ele) {
      label.appendChild(ele);
    }

    queriesFormEle?.appendChild(container);
  });
}

renderQueriesForm();

const defaultQueries = {
  size: "32",
  format: "png",
  rotate: "32",
};

type Query = {
  size: string;
  format: string;
  rotate: string;
};

let queries = { ...defaultQueries };

function updateQueryString(newQueries: Partial<Query>) {
  queries = { ...queries, ...newQueries };
  const queryString = Object.keys(queries)
    .map((key) => {
      const value = queries[key as keyof Query] as string;
      if (value === defaultQueries[key as keyof Query]) {
        return "";
      }
      return `${key}=${value}`;
    })
    .filter(Boolean)
    .join("&");

  const pureDomain = input.value?.split("?")[0] || DEFAULT_DOMAIN;
  updateInputValue(queryString ? `${pureDomain}?${queryString}` : pureDomain);
}

function updateInputValue(value: string) {
  input.value = value;
  handleUrlChange(value);
}
