<script setup lang="ts">
import { ref, computed, nextTick } from "vue";
import { useRoute } from "vue-router";

// Utility types
type Teacher = {
  created?: string;
  updated?: string;
  id: string;
  name: string;
  subjects: string[];
  schedule?: any;
  availabilities?: any;
};
const DAYS = ["Mon", "Tue", "Wed", "Thu", "Fri"];
// Table structure: lessons/details for a day
const SCHEDULE = [
  { type: "lesson", label: "Lesson 1", start: "09:00", end: "09:55" },
  { type: "lesson", label: "Lesson 2", start: "09:55", end: "10:50" },
  { type: "break", label: "Break", start: "10:50", end: "11:10" },
  { type: "lesson", label: "Lesson 3", start: "11:10", end: "12:05" },
  { type: "lesson", label: "Lesson 4", start: "12:05", end: "13:00" },
  { type: "lunch", label: "Lunch", start: "13:00", end: "13:40" },
  { type: "lesson", label: "Lesson 5", start: "13:40", end: "14:35" },
  { type: "lesson", label: "Lesson 6", start: "14:35", end: "15:30" },
];
// PocketBase client placeholder
const pb = usePocketbase();
const route = useRoute();

// Refs and state
const allTeachers = ref<Teacher[]>([]);
const selectedTeachers = ref<string[]>([]);
const activeTeacherId = ref<string | null>(null);
const activeTeacher = computed(
  () => allTeachers.value.find((t) => t.id === activeTeacherId.value) ?? null
);
const currentTimetable = ref<any>(null);

// Edit fields
const teacherName = ref("");
const subjectInput = ref("");
const teacherSubjects = ref<string[]>([]);

// Timetable (activated cells)
const timetableMatrix = ref<Record<string, boolean>>({});
function slotKey(dayIdx: number, periodIdx: number) {
  return `${dayIdx}-${periodIdx}`;
}

// Fetch teachers and timetable on mount!
onMounted(async () => {
  const teachers = await pb.collection("teachers").getFullList();
  allTeachers.value = teachers.map((t: any) => ({
    id: t.id,
    name: t.name,
    subjects: t.subjects || [],
    // FIX: here map .schedule to teacher.schedule not timetable!
    schedule: t.schedule || [],
    availabilities: t.availabilities || [],
  }));

  const timetableId = route.params.timetableId as string;
  if (timetableId) {
    const tt = await pb.collection("timetables").getOne(timetableId);
    currentTimetable.value = tt;
    // (rest unchanged)
    let teacherIds: string[] = [];
    if (typeof tt.teachers === "string") {
      try {
        teacherIds = JSON.parse(tt.teachers);
      } catch {
        teacherIds = [];
      }
    } else if (Array.isArray(tt.teachers)) {
      teacherIds = tt.teachers;
    }
    selectedTeachers.value = teacherIds.filter((id) =>
      allTeachers.value.some((t) => t.id === id)
    );
  }
});

// On teacher select: load info
function selectTeacher(teacher: Teacher) {
  activeTeacherId.value = teacher.id;
  teacherName.value = teacher.name;
  teacherSubjects.value = Array.from(teacher.subjects);
  console.log("selected ", activeTeacher.value);

  timetableMatrix.value = {}; // clear current matrix
  // FIX: use .schedule (which is correctly populated on allTeachers)
  if (teacher.schedule && Array.isArray(teacher.schedule)) {
    for (const slot of teacher.schedule) {
      timetableMatrix.value[slot] = true;
    }
  }
}

// Toggle a time slot
function toggleSlot(dayIdx: number, periodIdx: number) {
  const key = slotKey(dayIdx, periodIdx);
  timetableMatrix.value[key] = !timetableMatrix.value[key];
}

type AvailabilityRaw = {
  dow: number | string;
  start: string;
  end: string;
};
function groupAvailabilitiesToRanges(availSlots: string[]): AvailabilityRaw[] {
  // availSlots are ["0-2", ...]
  const byDay: Record<number, number[]> = {};
  // Group period indices by day
  for (const slot of availSlots) {
    const [d, p] = slot.split("-").map(Number);
    if (!byDay[d]) byDay[d] = [];
    byDay[d].push(p);
  }

  // For each day, create consecutive blocks
  const result: AvailabilityRaw[] = [];
  for (const dowStr in byDay) {
    const dow = Number(dowStr);
    const periods = byDay[dow].sort((a, b) => a - b);
    let rangeStart: number | null = null;
    let rangeEnd: number | null = null;

    for (let i = 0; i <= periods.length; i++) {
      const curr = periods[i];
      const prev = periods[i - 1];

      if (rangeStart === null) {
        rangeStart = curr;
        rangeEnd = curr;
      } else if (curr === prev + 1) {
        // Continue the range
        rangeEnd = curr;
      } else {
        // Close range and push
        if (rangeStart !== null && rangeEnd !== null) {
          const startTime = SCHEDULE[rangeStart].start;
          // End time is *end* of last period in block
          const endTime = SCHEDULE[rangeEnd].end;

          // In ISO8601 format (hour:minute:second, day only matters for date not time-of-day)
          result.push({
            dow: (dow + 1) % 7, // 0=Mon, 1=Tue, etc (assuming DAYS is Mon-Fri starting at 0)
            start: `${startTime}`, // e.g. "T09:00:00"
            end: `${endTime}`, // e.g. "T09:55:00"
          });
        }
        if (curr !== undefined) {
          rangeStart = curr;
          rangeEnd = curr;
        }
      }
    }
  }
  return result;
}

// Save timetable and availabilities
async function saveTeacherTimetable() {
  if (!activeTeacherId.value) return;
  const teacherIdx = allTeachers.value.findIndex(
    (t) => t.id === activeTeacherId.value
  );
  if (teacherIdx < 0) return;

  // Timetable slots ("dayIdx-periodIdx")
  const timetable = Object.keys(timetableMatrix.value).filter(
    (k) => timetableMatrix.value[k]
  );

  // Calculate available slots, i.e. all slots not in timetable
  const allSlotKeys: string[] = [];
  for (let d = 0; d < DAYS.length; d++) {
    for (let p = 0; p < SCHEDULE.length; p++) {
      allSlotKeys.push(slotKey(d, p));
    }
  }
  const availabilities = allSlotKeys.filter(
    (slot) => !timetableMatrix.value[slot]
  );
  // ---- AVAILABILITY RANGE GROUPING (the main part!) ----

  // Now, just before you save to DB:
  const availabilityRanges = groupAvailabilitiesToRanges(availabilities);

  // availabilityRanges is Array<AvailabilityRaw>
  // --- END ---

  // Save to DB (update in-place for UI)
  const updatedTeacher = {
    ...allTeachers.value[teacherIdx],
    name: teacherName.value,
    subjects: teacherSubjects.value,
    timetable,
    availabilities: availabilityRanges, // Still the slot IDs, optional for backwards compat
  };
  await pb.collection("teachers").update(activeTeacherId.value, {
    name: teacherName.value,
    subjects: teacherSubjects.value,
    schedule: timetable,
    availabilities: availabilityRanges, // or substitute your field name
  });
  allTeachers.value[teacherIdx] = updatedTeacher;
  alert("Saved!");
}

// Add/remove teacher to selection list for timetable
function selectTeacherForTimetable(id: string) {
  console.log(id);
  if (selectedTeachers.value.includes(id)) {
    selectedTeachers.value = selectedTeachers.value.filter((tid) => tid !== id);
  } else {
    selectedTeachers.value.push(id);
  }
}

// Save selected teachers to timetable (in timetableId param)
async function saveSelectedTeachersToTimetable() {
  const data = selectedTeachers.value.map((id) => {
    return id;
  });
  console.log(data);
  const timetableId = route.params.timetableId;
  await pb.collection("timetables").update(timetableId as string, {
    teachers: data,
  });
  alert("Assigned teachers to timetable.");
}

// Add new teacher
const newTeacherName = ref("");
const newTeacherSubjects = ref<string[]>([]);
const newTeacherSubjectInput = ref("");

async function addNewTeacher() {
  const teacher = await pb.collection("teachers").create({
    // name: newTeacherName.value,
    // subjects: newTeacherSubjects.value,
    name: "New Teacher",
    subjects: [],
  });
  // @ts-ignore
  allTeachers.value.push({ ...teacher, schedule: [], availabilities: [] });
  newTeacherName.value = "";
  newTeacherSubjects.value = [];
}

function onSubjectInputKeydown(e: KeyboardEvent) {
  if (e.key === "Enter" && subjectInput.value.trim()) {
    if (!teacherSubjects.value.includes(subjectInput.value.trim())) {
      teacherSubjects.value.push(subjectInput.value.trim());
    }
    subjectInput.value = "";
    nextTick(() => {
      (e.target as HTMLInputElement).focus();
    });
  }
}
function removeSubject(idx: number) {
  teacherSubjects.value.splice(idx, 1);
}

// Fetch teachers on mount
(async function fetchTeachers() {
  const teachers = await pb.collection("teachers").getFullList();
  allTeachers.value = teachers.map((t: any) => ({
    id: t.id,
    name: t.name,
    subjects: t.subjects || [],
    timetable: t.schedule || [],
    availabilities: t.availabilities || [],
  }));
})();
</script>
<template>
  <div class="flex flex-col gap-5 mt-10">
    <div class="flex gap-2">
      <NuxtLink to="/timetables"
        ><Icon class="size-6" name="material-symbols:arrow-left-alt-rounded"
      /></NuxtLink>
      <h2>Manage Teachers</h2>
    </div>
    <div class="flex flex-col gap-6 md:flex-row">
      <!-- Teacher List -->
      <Card class="flex flex-col justify-between p-4 w-72">
        <div>
          <h3 class="mb-2 text-lg font-semibold">All Teachers</h3>
          <Button class="w-full mb-3" @click="saveSelectedTeachersToTimetable">
            Save Teachers to Timetable
          </Button>
          <ul class="flex flex-col gap-2 mb-5">
            <li
              v-for="t in allTeachers"
              :key="t.id"
              :class="[
                'py-2 flex items-center justify-between cursor-pointer text-sm hover:bg-blue-100 transition border border-transparent p-2 rounded',
                activeTeacherId === t.id ? 'border-primary' : '',
              ]"
            >
              <div @click="selectTeacher(t)">
                {{ t.name }}
              </div>
              <input
                type="checkbox"
                :checked="selectedTeachers.includes(t.id)"
                @change="selectTeacherForTimetable(t.id)"
              />
            </li>
          </ul>
        </div>
        <div class="mt-8">
          <Button variant="outline" class="w-full" @click="addNewTeacher">
            Create Teacher
          </Button>
        </div>
      </Card>
      <!-- Timetable Editor -->
      <div class="md:w-3/4 rounded shadow px-6 py-4 min-h-[500px]">
        <div v-if="activeTeacher" class="mb-6">
          <div class="flex flex-col items-start gap-4 mb-4 md:flex-row">
            <div>
              <label class="block font-semibold">Teacher Name</label>
              <Input type="text" v-model="teacherName" />
            </div>
            <div class="flex-1">
              <label class="block font-semibold">Subjects</label>
              <div class="flex flex-wrap items-center gap-2">
                <Input
                  type="text"
                  placeholder="Add subject"
                  v-model="subjectInput"
                  @keydown="onSubjectInputKeydown"
                />
                <span
                  @click="removeSubject(i)"
                  v-for="(subject, i) in teacherSubjects"
                  :key="subject"
                  class="flex items-center gap-1 px-2 py-1 text-xs font-medium uppercase rounded-full cursor-pointer bg-primary"
                >
                  <p>{{ subject }}</p>
                </span>
              </div>
            </div>
          </div>
          <div>
            <h4 class="mb-2 font-bold">Timetable</h4>
            <div class="overflow-x-auto">
              <table
                class="min-w-[700px] table-auto border-collapse border-primary rounded shadow"
              >
                <thead>
                  <tr>
                    <th class="p-2 bg-gray-200 border border-primary"></th>
                    <th
                      v-for="(day, d) in DAYS"
                      :key="d"
                      class="p-2 bg-gray-200 border border-primary"
                    >
                      {{ day }}
                    </th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(period, pIdx) in SCHEDULE" :key="pIdx">
                    <td
                      :class="
                        period.type === 'break' || period.type === 'lunch'
                          ? 'bg-muted'
                          : ''
                      "
                      class="p-2 font-semibold bg-gray-100 border border-primary whitespace-nowrap"
                    >
                      <span
                        :class="
                          period.type === 'lesson'
                            ? 'text-blue-700'
                            : period.type === 'break'
                            ? ''
                            : 'text-green-700'
                        "
                      >
                        {{ period.label }}
                      </span>
                      <div class="text-xs text-gray-500">
                        {{ period.start }} - {{ period.end }}
                      </div>
                    </td>
                    <td
                      v-for="dayIdx in [0, 1, 2, 3, 4]"
                      :key="dayIdx"
                      :class="
                        period.type === 'break' || period.type === 'lunch'
                          ? 'bg-muted'
                          : ''
                      "
                      class="relative p-2 text-center border border-primary"
                    >
                      <button
                        :class="[
                          'w-8 h-8 rounded-full flex items-center justify-center mx-auto transition',
                          timetableMatrix[slotKey(dayIdx, pIdx)]
                            ? 'bg-blue-500 text-foreground bg-primary shadow-lg'
                            : 'bg-transparent border text-foreground hover:bg-primary hover:opacity-50 border-primary',
                        ]"
                        @click="toggleSlot(dayIdx, pIdx)"
                        title="Click to toggle lesson"
                      >
                        <span>{{
                          timetableMatrix[slotKey(dayIdx, pIdx)] ? "✓" : ""
                        }}</span>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <Button class="mt-5" @click="saveTeacherTimetable">
              Save Timetable & Availability
            </Button>
          </div>
        </div>
        <p v-else class="mt-8 italic text-gray-400">
          Select a teacher to edit their timetable and info.
        </p>
      </div>
    </div>
  </div>
</template>
