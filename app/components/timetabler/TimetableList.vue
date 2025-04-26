<template>
  <div>
    <ul class="flex gap-4">
      <Card
        v-for="t in timetables"
        :key="t.id"
        class="flex flex-col gap-2 p-4 w-72 bg-background"
      >
        <h3 class="font-bold">{{ t.name }}</h3>
        <div class="flex items-baseline justify-between">
          Teachers
          <Button
            @click="navigateTo(`/timetables/${t.id}/teachers`)"
            variant="ghost"
            >Edit</Button
          >
        </div>
        <div class="flex items-baseline justify-between">
          Exams to plan
          <Button
            @click="navigateTo(`/timetables/${t.id}/exams`)"
            variant="ghost"
            >Edit</Button
          >
        </div>
        <div v-if="t.timetable.length > 0" class="flex gap-2">
          <Button class="w-full">View Timetable</Button>
          <Button variant="outline"
            ><Icon class="size-5" name="material-symbols:rotate-left-rounded"
          /></Button>
          <Button variant="destructive" @click="deleteTimetable(t.id)"
            ><Icon name="material-symbols:delete-outline"
          /></Button>
        </div>
        <div class="flex w-full gap-2" v-else>
          <Button variant="secondary" class="w-full">Calculate</Button>
          <Button variant="destructive" @click="deleteTimetable(t.id)"
            ><Icon name="material-symbols:delete-outline"
          /></Button>
        </div>
      </Card>
    </ul>
  </div>
</template>
<script setup lang="ts">
const pb = usePocketbase();

interface Props {
  newTimetable: Timetable;
}
const props = defineProps<Props>();

export type Timetable = {
  id: string;
  name: string;
  exam_timetable: [];
  teachers: [];
  timetable: [];
};

watch(
  () => props.newTimetable,
  (curr, old) => {
    console.log("watch: new time", curr, old);
    if (!old && !curr) {
      return;
    }
    if (curr.id !== old.id) {
      timetables.value.unshift(curr);
    }
  }
);

async function getTimetables() {
  try {
    const records = await pb
      .collection("timetables")
      .getFullList<Timetable>(1, {
        sort: "-created",
      });
    timetables.value = records;
    console.log(records);
  } catch (err) {
    console.log(err);
  }
}

async function deleteTimetable(id: string) {
  try {
    await pb.collection("timetables").delete(id);
  } catch (err) {
    console.log(err);
  } finally {
    getTimetables();
  }
}

onMounted(async () => {
  await getTimetables();
});

const timetables = ref();

const tbls = ref([
  {
    id: "some",
    name: "Exam Empty",
    exam_timetable: [],
    teachers: [],
    timetable: [],
  },
  {
    id: "two",
    name: "Exam Complete",
    exam_timetable: [{ exam_board: "board" }],
    teachers: [],
    timetable: [{ name: "some" }],
  },
]);
</script>
