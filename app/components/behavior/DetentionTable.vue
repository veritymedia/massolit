<script setup lang="ts">
import type {
  ColumnFiltersState,
  ExpandedState,
  SortingState,
  VisibilityState,
} from "@tanstack/vue-table";
import { Button } from "@/components/ui/button";
import IconWrapper from "@/components/IconWrapper.vue";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import DOMPurify from "dompurify";

import { cn, valueUpdater } from "@/lib/utils";
import { CaretSortIcon, ChevronDownIcon } from "@radix-icons/vue";
import {
  createColumnHelper,
  FlexRender,
  getCoreRowModel,
  getExpandedRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useVueTable,
} from "@tanstack/vue-table";
import { type RecordListQueryParams, Record } from "pocketbase";
import { h, ref } from "vue";

export type BehaviorNote = {
  id: string;
  student_id: string;
  first_name: string;
  last_name: string;
  email: string;
  grade: string;
  incident_time: string;
  behavior_type: string;
  notes: string;
  next_step: string;
  next_step_date: string;
  author_id: number;
  reported_by: string;
  homeroom_advisor: string;
  visible_to_parents: boolean;
  visible_to_students: boolean;
  created_at: string;
  updated_at: string;
  action_complete: boolean;
};

const data2 = ref<BehaviorNote[]>([
  {
    id: "1",
    student_id: "S12345",
    first_name: "John",
    last_name: "Doe",
    email: "john.doe@example.com",
    grade: "10",
    incident_time: "2023-12-01T10:15:00Z",
    behavior_type: "Disrespect",
    notes: "Student was talking back to the teacher during class.",
    next_step: "Parent meeting scheduled",
    next_step_date: "2023-12-05",
    author_id: 101,
    reported_by: "Ms. Smith",
    homeroom_advisor: "Mr. Johnson",
    visible_to_parents: true,
    visible_to_students: false,
    created_at: "2023-12-01T10:30:00Z",
    updated_at: "2023-12-01T11:00:00Z",
    action_complete: false,
  },
  {
    id: "2",
    student_id: "S67890",
    first_name: "Emily",
    last_name: "Clark",
    email: "emily.clark@example.com",
    grade: "9",
    incident_time: "2023-11-28T14:45:00Z",
    behavior_type: "Cheating",
    notes: "Student was caught copying during the math exam.",
    next_step: "Detention assigned",
    next_step_date: "2023-11-30",
    author_id: 102,
    reported_by: "Mr. Brown",
    homeroom_advisor: "Ms. Green",
    visible_to_parents: true,
    visible_to_students: true,
    created_at: "2023-11-28T15:00:00Z",
    updated_at: "2023-11-28T15:30:00Z",
    action_complete: true,
  },
  {
    id: "3",
    student_id: "S54321",
    first_name: "Michael",
    last_name: "Lee",
    email: "michael.lee@example.com",
    grade: "11",
    incident_time: "2023-12-02T09:00:00Z",
    behavior_type: "Bullying",
    notes: "Student was observed intimidating another student in the hallway.",
    next_step: "Counselor session scheduled",
    next_step_date: "2023-12-06",
    author_id: 103,
    reported_by: "Mrs. White",
    homeroom_advisor: "Ms. Carter",
    visible_to_parents: false,
    visible_to_students: false,
    created_at: "2023-12-02T09:30:00Z",
    updated_at: "2023-12-02T10:00:00Z",
    action_complete: false,
  },
]);

const data = ref<BehaviorNote[]>([]);

const toggleComplete = async (id: string) => {
  console.log(`Row with id ${id} set to complete.`);

  const row = data.value.find((row) => row.id === id);
  if (row) {
    if (row.action_complete === true) {
      const res = await setActionCompleteOnRecord(id, false);
      if (res) {
        row.action_complete = res.action_complete;
      }
    } else {
      const res = await setActionCompleteOnRecord(id, true);
      if (res) {
        row.action_complete = res.action_complete;
      }
    }
    console.log("Note ID: ", row.id);
  }
};
const columnHelper = createColumnHelper<BehaviorNote>();

const columns = [
  columnHelper.accessor("action_complete", {
    header: ({ table }) =>
      h(IconWrapper, {
        name: "material-symbols:check-rounded",
        class: "w-5 h-5",
      }),
    cell: ({ row }) => {
      return h(
        "div",
        {
          class: "",
        },
        row.original["action_complete"]
          ? h(IconWrapper, {
              name: "material-symbols:check-rounded",
              class: "w-5 h-5 opacity-[50%]",
            })
          : h(IconWrapper, {
              name: "material-symbols:warning-rounded",
              class: "w-5 h-5",
            }),
      );
    },
    enableSorting: true,
    enableHiding: false,
  }),
  columnHelper.accessor("grade", {
    header: ({ table }) => {
      return h("div", {}, "Year");
    },
  }),
  columnHelper.display({
    id: "full_name",
    header: () => h("div", { class: " left-right" }, "Full Name"),
    cell: ({ row }) => {
      return h(
        "div",
        { class: "text-left font-medium" },
        row.original["first_name"] + " " + row.original["last_name"],
      );
    },
  }),
  columnHelper.accessor("incident_time", {
    header: ({ table }) => {
      return h("div", {}, "Date");
    },
    cell: ({ row }) => {
      return h(
        "p",
        {},
        new Intl.DateTimeFormat("en-UK", {
          day: "numeric",
          month: "numeric",
          year: "2-digit",
        }).format(new Date(row.getValue("incident_time"))),
      );
    },
  }),
  columnHelper.display({
    id: "set_complete",
    header: () => h("div", {}, "Actions"),
    cell: ({ row, getValue }) => {
      const buttonValue = () => {
        return row.original["action_complete"] ? "Undo" : "Done";
      };

      return h(
        Button,
        {
          onClick: () => {
            toggleComplete(row.original.id);
          },
          size: "xs",
          variant:
            row.original["action_complete"] === true ? "ghost" : "default",
          class: "text-xs",
        },
        buttonValue,
      );
    },
  }),
  columnHelper.accessor("notes", {
    header: ({ table }) => {
      return h("div", {}, "Note");
    },
    cell: ({ getValue }) => {
      // Sanitize HTML while keeping inline styles
      const cleanHtml = DOMPurify.sanitize(getValue(), {
        ALLOWED_TAGS: ["p", "span", "strong", "em", "u", "b", "i"],
        ALLOWED_ATTR: ["style", "class"],
      });

      return h("div", {
        innerHTML: cleanHtml,
        class: "rich-text-content", // Optional: add a class for additional styling
      });
    },
  }),
];

const sorting = ref<SortingState>([]);
const columnFilters = ref<ColumnFiltersState>([]);
const columnVisibility = ref<VisibilityState>({});
const rowSelection = ref({});
const expanded = ref<ExpandedState>({});

const table = useVueTable({
  data,
  columns,
  getCoreRowModel: getCoreRowModel(),
  getPaginationRowModel: getPaginationRowModel(),
  getSortedRowModel: getSortedRowModel(),
  getFilteredRowModel: getFilteredRowModel(),
  getExpandedRowModel: getExpandedRowModel(),
  onSortingChange: (updaterOrValue) => valueUpdater(updaterOrValue, sorting),
  onColumnFiltersChange: (updaterOrValue) =>
    valueUpdater(updaterOrValue, columnFilters),
  onColumnVisibilityChange: (updaterOrValue) =>
    valueUpdater(updaterOrValue, columnVisibility),
  onRowSelectionChange: (updaterOrValue) =>
    valueUpdater(updaterOrValue, rowSelection),
  onExpandedChange: (updaterOrValue) => valueUpdater(updaterOrValue, expanded),
  state: {
    get sorting() {
      return sorting.value;
    },
    get columnFilters() {
      return columnFilters.value;
    },
    get columnVisibility() {
      return columnVisibility.value;
    },
    get rowSelection() {
      return rowSelection.value;
    },
    get expanded() {
      return expanded.value;
    },
    columnPinning: {
      left: ["status"],
    },
  },
});

const pb = usePocketbase();

async function getBehaviourNotes(page: number = 1) {
  try {
    const params: RecordListQueryParams = {
      sort: "+action_complete,-next_step_date",
    };

    const res = await pb
      .collection("behavior_notes")
      .getList(page, 200, params);
    console.log(res);
    data.value = res.items;
  } catch (err) {
    console.log(err);
  }
}

async function setActionCompleteOnRecord(
  recordId: string,
  actionState: boolean,
): Promise<Record | undefined> {
  try {
    const res = await pb.collection("behavior_notes").update(recordId, {
      action_complete: actionState,
    });
    console.log("record updated: ", res);
    return res;
  } catch (err) {
    console.log(err);
    return undefined;
  }
}

onMounted(async () => {
  await getBehaviourNotes();
});
</script>

<template>
  <div class="w-full">
    <BehaviorDetentionTablePaginateControl :table="table" />

    <div class="border rounded-md text-xs">
      <Table class="border-[white]">
        <TableHeader>
          <TableRow
            v-for="headerGroup in table.getHeaderGroups()"
            :key="headerGroup.id"
          >
            <TableHead
              v-for="header in headerGroup.headers"
              :key="header.id"
              :data-pinned="header.column.getIsPinned()"
              :class="
                cn(
                  { 'sticky bg-background/95': header.column.getIsPinned() },
                  header.column.getIsPinned() === 'left' ? 'left-0' : 'right-0',
                )
              "
            >
              <FlexRender
                v-if="!header.isPlaceholder"
                :render="header.column.columnDef.header"
                :props="header.getContext()"
              />
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <template v-if="table.getRowModel().rows?.length">
            <template v-for="row in table.getRowModel().rows" :key="row.id">
              <TableRow
                class=""
                :data-state="row.getIsSelected() && 'selected'"
              >
                <TableCell
                  v-for="cell in row.getVisibleCells()"
                  :key="cell.id"
                  class=""
                  :data-pinned="cell.column.getIsPinned()"
                  :class="
                    cn(
                      { 'sticky bg-background/95': cell.column.getIsPinned() },
                      cell.column.getIsPinned() === 'left'
                        ? 'left-0'
                        : 'right-0',
                    )
                  "
                >
                  <FlexRender
                    :render="cell.column.columnDef.cell"
                    :props="cell.getContext()"
                  />
                </TableCell>
              </TableRow>
              <TableRow v-if="row.getIsExpanded()">
                <TableCell :colspan="row.getAllCells().length">
                  {{ JSON.stringify(row.original) }}
                </TableCell>
              </TableRow>
            </template>
          </template>

          <TableRow v-else>
            <TableCell :colspan="columns.length" class="h-24 text-center">
              No results.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
