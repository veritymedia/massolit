<template>
  <div class="flex flex-col gap-5">
    <h1>Create New Exam List</h1>
    <div class="flex flex-col gap-2">
      <div class="flex items-end gap-2 mb-5">
        <div class="flex flex-col w-auto">
          Exam Board
          <Input class="w-64" v-model="exam_board"></Input>
        </div>

        <div class="flex flex-col w-auto">
          Qualification
          <Input class="w-64" v-model="qualification"></Input>
        </div>

        <div class="flex flex-col w-auto">
          Session
          <Input class="w-64" v-model="session"></Input>
        </div>
        <div
          v-if="error"
          class="p-2 rounded h-min bg-destructive text-destructive-foreground"
        >
          {{ error }}
        </div>
      </div>
    </div>

    <h2>Upload CSV File</h2>
    <div
      v-if="csvData.length === 0"
      class="p-6 mb-6 text-center border-2 border-gray-400 border-dashed rounded-lg"
      @dragover.prevent
      @drop.prevent="handleFileDrop"
    >
      <input
        type="file"
        ref="fileInput"
        accept=".csv"
        class="hidden"
        @change="handleFileSelect"
      />
      <Button @click="triggerFileInput"> Select CSV File </Button>
      <p class="mt-2 text-gray-500">or drag and drop your file here</p>
      <p v-if="fileName" class="mt-2 text-sm text-gray-700">
        Selected file: {{ fileName }}
      </p>
    </div>

    <div class="flex items-baseline mb-3 space-x-4">
      <Input
        v-model="searchTerm"
        type="text"
        placeholder="Search subject..."
        class="px-3 py-1 text-sm border border-gray-300 rounded-md w-72"
      />
      <Button
        variant="destructive"
        class="disabled:bg-gray-400"
        :disabled="!anySelected"
        @click="multiDelete"
      >
        Delete Selected ({{ selectedRowIds.size }})
      </Button>
    </div>
    <div v-if="parsedData.length > -1" class="mb-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold">Exams ({{ parsedData.length }})</h3>
        <div class="flex items-center gap-2">
          <Button variant="secondary" @click="startAddNewRow"
            >Create New Exam</Button
          >
          <Button @click="emitParsedContent">Save Exam List</Button>
        </div>
      </div>

      <div
        class="overflow-x-auto relative shadow-2xl max-h-[60vh] text-sm bg-background p-4 rounded-2xl"
      >
        <table class="min-w-full">
          <thead>
            <tr>
              <th class="px-2 py-2 text-left border-b">
                <input
                  type="checkbox"
                  :checked="allVisibleSelected"
                  @change="toggleSelectAllVisibleRows($event.target.checked)"
                />
              </th>
              <th class="px-2 py-0 text-left border-b">Subject</th>
              <th class="px-2 py-0 text-left border-b">Start Time</th>
              <th class="px-2 py-0 text-left border-b">Duration</th>
              <th class="px-2 py-0 text-left border-b">Room</th>
              <th class="px-2 py-0 text-left border-b">Exam Code</th>
              <th class="px-2 py-0 text-left border-b">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="newRow">
              <td class="border-b"></td>
              <td class="px-2 py-0 border-b">
                <Input v-model="newRow.subject" type="text" />
              </td>
              <td class="px-2 py-0 border-b">
                <Input v-model="newRow.start" type="datetime-local" />
              </td>
              <td class="px-2 py-0 border-b">
                <Input
                  v-model="newRow.duration"
                  type="text"
                  placeholder="01:30"
                />
              </td>
              <td class="px-2 py-0 border-b">
                <Input v-model="newRow.room" type="text" />
              </td>
              <td class="px-2 py-0 border-b">
                <Input v-model="newRow.examCode" type="text" />
              </td>
              <td class="flex gap-2 px-2 py-2 border-b">
                <Button size="xs" @click="saveNewRow"> Save </Button>
                <Button size="xs" variant="secondary" @click="cancelNewRow"
                  >Cancel</Button
                >
              </td>
            </tr>

            <tr
              v-for="item in filteredRows"
              :key="item.id"
              class="hover:bg-zinc-900"
            >
              <template v-if="editRowId === item.id">
                <td class="px-2 py-1 border-b">
                  <Input
                    type="checkbox"
                    :checked="selectedRowIds.has(item.id)"
                    @change="toggleSelectRow(item.id, $event.target.checked)"
                  />
                </td>
                <td class="px-2 py-0 border-b">
                  <Input v-model="editBuffer.subject" type="text" />
                </td>
                <td class="px-2 py-0 border-b">
                  <Input v-model="editBuffer.start" type="datetime-local" />
                </td>
                <td class="px-2 py-0 border-b">
                  <Input
                    v-model="editBuffer.duration"
                    type="text"
                    placeholder="01:30"
                  />
                </td>
                <td class="px-2 py-0 border-b">
                  <Input v-model="editBuffer.room" type="text" />
                </td>
                <td class="px-2 py-0 border-b">
                  <Input v-model="editBuffer.examCode" type="text" />
                </td>
                <td class="flex gap-2 px-2 py-2 border-b">
                  <Button size="xs" @click="saveEditRow"> Save </Button>
                  <Button size="xs" variant="secondary" @click="cancelEditRow"
                    >Cancel</Button
                  >
                </td>
              </template>
              <template v-else>
                <td class="px-2 py-2 border-b">
                  <input
                    type="checkbox"
                    :checked="selectedRowIds.has(item.id)"
                    @change="toggleSelectRow(item.id, $event.target.checked)"
                  />
                </td>
                <td class="px-2 py-0 border-b">{{ item.subject }}</td>
                <td class="px-2 py-0 border-b">
                  {{ new Date(item.start).toLocaleString() }}
                </td>
                <td class="px-2 py-0 border-b">{{ item.duration }}</td>
                <td class="px-2 py-0 border-b">{{ item.room }}</td>
                <td class="px-2 py-0 border-b">{{ item.examCode || "—" }}</td>
                <td class="px-2 py-0 border-b">
                  <div class="flex items-center space-x-2">
                    <Button
                      @click="startEditRow(item)"
                      size="xs"
                      variant="outline"
                      >Edit</Button
                    >
                    <Button
                      variant="outline"
                      size="xs"
                      @click="duplicateExam(item)"
                      >Duplicate</Button
                    >
                    <Button
                      variant="destructive"
                      size="xs"
                      @click="deleteExam(item.id)"
                      >Delete</Button
                    >
                  </div>
                </td>
              </template>
            </tr>
          </tbody>
        </table>
        <div v-if="filteredRows.length === 0" class="mt-3 text-gray-500">
          No matching records.
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Imports
import { ref, computed } from "vue";
import Papa from "papaparse";
import Fuse from "fuse.js";
import { shapeData } from "./csv-utils";

// Emits
const emit = defineEmits(["exams:add"]);

// The new refs for user inputs
const exam_board = ref("Pearson");
const qualification = ref("iGCSE");
const session = ref("Summer 2025");

const allowedQualifications = ["iGCSE", "GCSE", "iAL", "AL"];

function emitParsedContent() {
  // Clear errors before checking
  error.value = "";

  // Trimmed values
  const ex = exam_board.value.trim();
  const qu = qualification.value.trim();
  const se = session.value.trim();

  if (!ex || !qu || !se) {
    error.value = "Please fill in all fields.";
    return;
  }
  if (!allowedQualifications.includes(qu)) {
    error.value = "Qualification must be one of: iGCSE, GCSE, iAL, AL";
    return;
  }

  emit("exams:add", {
    data: parsedData.value,
    exam_board: ex,
    qualification: qu,
    session: se,
  });
  // Optionally clear error or form fields here
  csvData.value = [];
  parsedData.value = [];
  fileName.value = "";
  error.value = "";
  validationResults.value = [];
}
// Props
const props = defineProps({
  config: {
    type: Object,
    default: () => ({
      morning: "10:00",
      afternoon: "12:00",
    }),
  },
});

// File handling
const fileInput = ref<HTMLInputElement | null>(null);
const fileName = ref("");
const csvData = ref<any[]>([]);
const error = ref<string>("");
const validationResults = ref<any[]>([]);
const parsedData = ref<any[]>([]);

// Extra: show/hide headers warning if needed

const expectedHeaders = [
  "date",
  // "Exam series",
  // "Board",
  // "Qual",
  "examCode",
  "subject",
  // "Title",
  "time",
  "duration",
];

function triggerFileInput() {
  fileInput.value?.click();
}
function handleFileSelect(event: Event) {
  const files = (event.target as HTMLInputElement).files;
  if (files && files.length > 0) {
    processFile(files[0]);
  }
}
function handleFileDrop(event: DragEvent) {
  if (event.dataTransfer && event.dataTransfer.files.length > 0) {
    processFile(event.dataTransfer.files[0]);
  }
}
function processFile(file: File) {
  if (file.type !== "text/csv" && !file.name.endsWith(".csv")) {
    error.value = "Please upload a CSV file.";
    fileName.value = "";
    return;
  }
  fileName.value = file.name;
  error.value = "";
  validationResults.value = [];
  parsedData.value = [];

  Papa.parse(file, {
    header: true,
    skipEmptyLines: true,
    complete: (results) => {
      csvData.value = results.data;
      validateHeaders(results.meta.fields);
      if (!error.value) {
        parseToExamRawFormat();
      }
    },
    error: (err) => {
      error.value = `Error parsing CSV: ${err.message}`;
    },
  });
}

function validateHeaders(headers?: string[]) {
  validationResults.value = [];
  if (!headers || headers.length === 0) {
    error.value = "No headers found in CSV file.";
    return;
  }

  // Check for all expected headers
  expectedHeaders.forEach((header) => {
    const found = headers.includes(header);
    validationResults.value.push({
      valid: found,
      message: found
        ? `✓ Found header: ${header}`
        : `✗ Missing header: ${header}`,
    });
  });

  // Check for any unexpected headers
  headers.forEach((header) => {
    if (!expectedHeaders.includes(header)) {
      validationResults.value.push({
        valid: false,
        message: `! Unexpected header: ${header}`,
      });
    }
  });

  // Set error if any required headers are missing
  const missingHeaders = expectedHeaders.filter((h) => !headers.includes(h));
  if (missingHeaders.length > 0) {
    error.value =
      "CSV file is missing required headers: " + missingHeaders.join(", ");
  }
}

// ---- CSV Parsing/Transformation ----
function parseToExamRawFormat() {
  parsedData.value = csvData.value
    .map((row, index) => {
      return shapeData(row);
    })
    .filter((item) => item !== null);
}

// Editing/Adding Logic
const editRowId = ref<string | null>(null);
const editBuffer = ref<any>({});
const newRow = ref<any | null>(null);

function startEditRow(row: any) {
  editRowId.value = row.id;
  editBuffer.value = { ...row };
}
function cancelEditRow() {
  editRowId.value = null;
  editBuffer.value = {};
}
function saveEditRow() {
  if (!editRowId.value) return;
  const idx = parsedData.value.findIndex((e) => e.id === editRowId.value);
  if (idx !== -1) {
    parsedData.value[idx] = { ...editBuffer.value };
  }
  cancelEditRow();
}

function startAddNewRow() {
  const now = new Date();
  newRow.value = {
    id: `exam-${Date.now()}`,
    subject: "",
    start: now.toISOString().slice(0, 16),
    duration: "",
    room: "",
    examCode: "",
  };
}
function cancelNewRow() {
  newRow.value = null;
}
function saveNewRow() {
  if (!newRow.value.subject || !newRow.value.start || !newRow.value.duration) {
    alert("Please fill all required fields");
    return;
  }
  parsedData.value.push({ ...newRow.value });
  cancelNewRow();
}

function duplicateExam(row: any) {
  newRow.value = {
    ...row,
    id: `exam-${Date.now()}`,
  };
  window.scrollTo({ top: 0, left: 0, behavior: "smooth" });
}
function deleteExam(examId: string) {
  const idx = parsedData.value.findIndex((ex) => ex.id === examId);
  if (idx !== -1) {
    parsedData.value.splice(idx, 1);
    if (editRowId.value === examId) cancelEditRow();
  }
}

// --- Filtering and Fuzzy Search ---
const searchTerm = ref("");
const selectedRowIds = ref<Set<string>>(new Set());

const filteredRows = computed(() => {
  if (!searchTerm.value.trim()) return parsedData.value;
  const fuse = new Fuse(parsedData.value, {
    keys: ["subject"],
    threshold: 0.2,
  });
  return fuse.search(searchTerm.value.trim()).map((r: any) => r.item);
});

// Selection
const allVisibleSelected = computed(() => {
  if (filteredRows.value.length === 0) return false;
  return filteredRows.value.every((row) => selectedRowIds.value.has(row.id));
});
const anySelected = computed(() => selectedRowIds.value.size > 0);

function toggleSelectRow(id: string, checked: boolean) {
  if (checked) {
    selectedRowIds.value = new Set(selectedRowIds.value).add(id);
  } else {
    const s = new Set(selectedRowIds.value);
    s.delete(id);
    selectedRowIds.value = s;
  }
}
function toggleSelectAllVisibleRows(checked: boolean) {
  let s = new Set(selectedRowIds.value);
  if (checked) {
    filteredRows.value.forEach((row) => s.add(row.id));
  } else {
    filteredRows.value.forEach((row) => s.delete(row.id));
  }
  selectedRowIds.value = s;
}
function multiDelete() {
  parsedData.value = parsedData.value.filter(
    (row) => !selectedRowIds.value.has(row.id),
  );
  selectedRowIds.value = new Set();
  searchTerm.value = "";
}
</script>
