<template>
  <div class="flex flex-col gap-5 mt-10">
    <div class="flex gap-2">
      <NuxtLink to="/timetables"
        ><Icon class="size-6" name="material-symbols:arrow-left-alt-rounded"
      /></NuxtLink>
      <h2>Exams</h2>
    </div>

    <div v-if="timetable" class="flex gap-2">
      <Card
        v-for="et in timetable.expand.exam_timetables"
        :key="et.id"
        class="flex flex-col w-auto gap-2 p-4"
      >
        <div>
          <div>{{ et.exam_board }} {{ et.qualification }} {{ et.session }}</div>
          <div>Exams: {{ et.data.length }}</div>
        </div>
        <div class="flex gap-2">
          <Button
            v-if="et.id !== $route.params.timetableListId"
            variant="outline"
            class="w-full"
            @click="
              navigateTo(
                `/timetables/${$route.params.timetableId}/exams/${et.id}`,
              )
            "
            >Edit</Button
          >
          <Button
            v-else
            class="w-full"
            @click="navigateTo(`/timetables/${$route.params.timetableId}`)"
            >Cancel
          </Button>

          <Button variant="destructive" @click="deleteExamTimetable(et.id)"
            ><Icon name="material-symbols:delete"
          /></Button>
        </div>
      </Card>
    </div>

    <TimetablerExamUpload
      v-if="!$route.params.timetableListId"
      @exams:add="addExams"
    />
    <NuxtPage v-if="$route.params.timetableListId" />
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
      { expand: "exam_timetables" },
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
