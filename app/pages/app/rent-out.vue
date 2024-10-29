<template>
  <div class="flex flex-col items-center">
    <Icon name="material-symbols:check-circle" class="w-12 h-12 opacity-70" />
    <p class="w-full text-center my-5 text-lg font-bold uppercase">
      Book available
    </p>
    <BookCard
      :book="{
        id: props.rentedBookStatus.bookId,
        isbn: props.rentedBookStatus.book?.isbn,
        title: props.rentedBookStatus.book?.title,
        cover_url: props.rentedBookStatus.book?.cover_url,
      }"
    />
  </div>
  <div class="w-full mt-5">
    <Dialog v-model:open="bookRentModal" class="">
      <DialogTrigger class="sticky bottom-5" as-child>
        <Button class="w-full" variant="secondary">Select Student</Button>
      </DialogTrigger>
      <DialogContent class="h-full flex flex-col items-center justify-between">
        <DialogHeader>
          <DialogTitle>Rent out book</DialogTitle>
        </DialogHeader>

        <div class="w-full flex flex-col items-center gap-4">
          <div class="w-full">
            <div
              class="w-full flex items-center items-evenly border border-muted rounded-lg"
            >
              <Input
                @input="handleInput"
                class="border-0"
                placeholder="Search for student"
                v-model="studentSearchTerm"
              />
              <Icon v-if="isLoading" name="line-md:loading-loop" class="mx-2" />
            </div>
            <Card
              v-if="studentList.length !== 0 && studentSearchTerm.length !== 0"
              class="mt-2"
            >
              <ul class="flex flex-col">
                <li
                  class="hover:bg-muted cursor-pointer first:mt-3 last:mb-3 py-0.5 px-3"
                  @click="handleStudentSelect(student)"
                  v-for="student in studentList"
                  :key="student.id"
                >
                  {{ student.first_name }} {{ student.last_name }}
                </li>
              </ul>
            </Card>
          </div>
          <h2 class="text-lg font-bold"></h2>
          <p class="text-center">
            The book
            <span class="font-bold">
              {{ props.rentedBookStatus.book?.title }}
            </span>
            with code
            <span class="font-bold bg-primary px-1 rounded-full">
              {{ props.rentedBookStatus.bookId }}</span
            >
            will be assigned to:
          </p>
          <div class="text-center" v-if="selectedStudent && selectedStudent.id">
            <span class="font-bold uppercase">
              {{ selectedStudent.first_name }} {{ selectedStudent.last_name }}
            </span>
            <br />
            ({{ selectedStudent.class_grade }})
          </div>

          <div>
            <Button variant="ghost" @click="closeDialog">Cancel</Button>
            <Button :disabled="!selectedStudent" @click="handleBookLease"
              >Register Book</Button
            >
          </div>
        </div>
        <DialogFooter class="w-full flex items-center"> </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import type { Record } from "pocketbase";
import type { RentedBookStatus as BookStatus } from "../app.vue";

interface Props {
  rentedBookStatus: BookStatus;
}

const isLoading = ref(false);
const managebacResult = ref();

const studentList = computed(() => {
  if (!managebacResult.value) {
    return [];
  }
  return managebacResult.value;
});

const selectedStudent = ref();
function handleStudentSelect(student: any) {
  console.log("Selected student: ", student.email);
  selectedStudent.value = student;
  studentSearchTerm.value = "";
  managebacResult.value = {};
}
const studentSearchTerm = ref("");

async function searchManagebacStudents() {
  if (studentSearchTerm.value.length < 3) {
    console.log("Search term too short. Must be 3+");
    return;
  }
  isLoading.value = true;
  try {
    const res = await pb.send(
      `/managebac/students?q=${studentSearchTerm.value}`,
      {},
    );
    console.log(res);
    managebacResult.value = res.students;
  } catch (err) {
    console.log(err);
  } finally {
    isLoading.value = false;
  }
}

let timeout = null;
const handleInput = () => {
  clearTimeout(timeout);
  timeout = setTimeout(async () => {
    // studentListOpen.value = true;
    console.log(studentSearchTerm.value);
    await searchManagebacStudents();
  }, 750);
};

async function handleBookLease() {
  try {
    const data = {
      rented_to: selectedStudent.value.id,
      book_instance: props.rentedBookStatus.book_instance?.id,
    };
    console.log("creating rental: ", data);
    const res = await pb.collection("rentals").create(data);
    console.log(res);
    if (res.id) {
      await navigateTo("/app");
    }
  } catch (err) {}
}

const props = withDefaults(defineProps<Props>(), {});

const pb = usePocketbase();

const { isOpen: bookRentModal, closeDialog, openDialog } = useDialogState();

onMounted(async () => {
  openDialog();
  // await searchManagebacStudents();
  // closeStudentSelectDialog();
});
</script>
