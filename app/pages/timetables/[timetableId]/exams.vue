<template>
  <div class="flex flex-col gap-5 mt-10 bg-white">
    <div class="flex gap-2">
      <NuxtLink to="/timetables"
        ><Icon class="size-6" name="material-symbols:arrow-left-alt-rounded"
      /></NuxtLink>
      <h2>Exams</h2>
    </div>

    <div v-if="timetable" class="flex gap-2">
      <Card
        v-for="t in timetable.expand.exam_timetables"
        :key="t.id"
        class="flex flex-col w-auto gap-2 p-4"
      >
        <div>
          <div>{{ t.exam_board }} {{ t.qualification }} {{ t.session }}</div>
          <div>Exams: {{ t.data.length }}</div>
        </div>

        <Button
          variant="destructive"
          size="sm"
          @click="deleteExamTimetable(t.id)"
          >DELETE</Button
        >
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
  try {
    console.log(data);

    const record = await pb.collection("exam_timetables").create(data);

    const tb = await pb
      .collection("timetables")
      .getOne(route.params.timetableId as string);

    tb.exam_timetables.push(record.id);

    const update = await pb.collection("timetables").update(
      route.params.timetableId as string,
      {
        ...tb,
      },
      { expand: "exam_timetables" }
    );

    timetable.value = update;
    console.log("Update: ", update);
  } catch (err) {
    console.log(err);
  }
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
  await pb.collection("exam_timetables").delete(id);
  await getTimetable();
}

onMounted(async () => {
  await getTimetable();
});
</script>
