<template>
  <div class="container mt-10 bg-white">
    <h2>Exams</h2>

    <div v-if="timetable" class="flex">
      <Card
        v-for="t in timetable.expand.exam_timetables"
        :key="t.id"
        class="w-auto p-2"
      >
        <div>
          {{ t.exam_board }} {{ t.qualification }} {{ t.session }} Exams:
          {{ t.data.length }}
        </div>
        <div><Button @click="deleteExamTimetable(t.id)">DELETE</Button></div>
      </Card>
    </div>

    <TimetablerExamUpload @exams:add="addExams" />
  </div>
</template>

<script lang="ts" setup>
const pb = usePocketbase();
const route = useRoute();
const timetable = ref();

async function addExams(data: any) {
  console.log(data);

  const record = await pb.collection("exam_timetables").create(data);

  const timetable = await pb
    .collection("timetables")
    .getOne(route.params.timetableId as string);

  timetable.exam_timetables.push(record.id);

  const update = await pb
    .collection("timetables")
    .update(route.params.timetableId as string, {
      ...timetable,
    });

  console.log(update);
}

async function getTimetable() {
  const records = await pb
    .collection("timetables")
    .getOne(route.params.timetableId as string, {
      expand: "exam_timetables",
    });
  console.log(records);
  timetable.value = records;
}

async function deleteExamTimetable(id: string) {
  pb.collection("exam_timetables").delete(id);
}

onMounted(async () => {
  await getTimetable();
});
</script>
