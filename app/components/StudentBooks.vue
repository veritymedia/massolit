<script lang="ts" setup>
const pb = usePocketbase();
const isLoading = ref(false);
const managebacResult = ref<ManagebacStudent[]>();
const selectedStudent = ref<ManagebacStudent>();
const studentSearchTerm = ref();
const studentRentalList = ref<RentalItem[] | undefined>(undefined);

type ManagebacStudent = {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  archived: boolean;
  ui_language: string;
  created_at: string;
  updated_at: string;
  student_id: string;
  identifier: string;
  last_accessed_at: string;
  gender: "Female" | "Male" | "Other" | "Prefer not to say";
  homeroom_advisor_id: number;
  class_grade: string;
  program: string;
  year_group_id: number;
  parent_ids: number[];
  street_address: string;
  city: string;
  state: string;
  zipcode: string;
  country: string;
  nationalities: string[];
  languages: string[];
  graduating_year: number;
};

const studentList = computed(() => {
  if (!managebacResult.value) {
    return [];
  }
  return managebacResult.value;
});

let timeout: any = null;

const handleInput = () => {
  clearTimeout(timeout);
  timeout = setTimeout(async () => {
    console.log(studentSearchTerm.value);
    await searchManagebacStudents();
  }, 750);
};

function handleStudentSelect(student: any) {
  console.log("Selected student: ", student.email);
  selectedStudent.value = student;
  studentSearchTerm.value = "";
  managebacResult.value = [];
}

async function searchManagebacStudents() {
  console.log("Search term: ", studentSearchTerm.value);

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

watch(selectedStudent, (v) => {
  if (v === undefined) {
    return;
  }
  getRentalsByStudentId(v.id);
});

type RentalItem = {
  book_instance: string;
  collectionId: string;
  collectionName: "rentals";
  created: string;
  expand: {
    book_instance?: BookInstance;
  };
  id: string;
  rented_to?: string;
  updated: string;
};

type BookInstance = {
  book: string;
  book_code: string;
  collectionId: string;
  collectionName: "book_instances";
  created: string;
  expand: {
    book?: Book;
  };
  id: string;
  updated: string;
};

type Book = {
  collectionId: string;
  collectionName: "books";
  cover_url: string;
  created: string;
  expand: {}; // No expanded properties in this example, but may add nested expansions if needed
  id: string;
  isbn: number;
  title: string;
  updated: string;
};

type RentalResponse = {
  items: RentalItem[];
  page: number;
  perPage: number;
  totalItems: number;
  totalPages: number;
};

async function getRentalsByStudentId(studentId: number) {
  try {
    const res: RentalResponse = await pb.collection("rentals").getList(1, 50, {
      filter: `rented_to="${studentId}"`,
      expand: "book_instance,book_instance.book",
    });
    console.log("Student books", res);
    studentRentalList.value = res.items;
  } catch (err) {
    console.log(err);
  }
}

const isDeleteMode = ref<boolean>(false);

const awaitingDelete = ref(false);

async function handleRentalReturn(id: string) {
  try {
    const deleted = await pb.collection("rentals").delete(id);
    if (deleted) {
      const i = studentRentalList.value?.findIndex((v) => {
        return v.id === id;
      });
      if (i !== -1 && i !== undefined) {
        studentRentalList.value?.splice(i, 1);
      }
    }
  } catch (err) {
    console.log(err);
  }
}

function dateToLocaleString(date: string): string {
  return new Date(date).toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}
</script>

<template>
  <div>
    <h2 class="text-lg font-bold">Search books by student</h2>
    <div class="w-full relative">
      <div
        class="w-full mt-2 flex items-center items-evenly border border-muted rounded-lg"
      >
        <Input
          @input="handleInput"
          class="border-0 w-full"
          placeholder="Search for student"
          v-model="studentSearchTerm"
        />
        <Icon v-if="isLoading" name="line-md:loading-loop" class="mx-2" />
      </div>
      <Card
        v-if="studentList.length !== 0 && studentSearchTerm.length !== 0"
        class="mt-2 absolute top-10 bg-background z-50 border-primary w-full overflow-hidden"
      >
        <ul class="flex flex-col">
          <li
            class="hover:bg-muted cursor-pointer first:mt-3 last:mb-3 py-0.5 px-3"
            @click="handleStudentSelect(student)"
            v-for="student in studentList"
            :key="student.id"
          >
            {{ student.first_name }} {{ student.last_name }}
            <span class="text-xs pl-2 uppercase font-semibold text-[gray]"
              >{{ student.class_grade }}
            </span>
          </li>
        </ul>
      </Card>
    </div>

    <div
      v-if="studentRentalList"
      class="w-full h-0.5 bg-accent mt-4 rounded-full"
    ></div>
    <div class="mt-5" v-if="studentRentalList !== undefined && selectedStudent">
      <div class="flex justify-between items-center mb-3">
        <p class="font-bold mb-1">
          {{ selectedStudent.first_name + " " + selectedStudent.last_name }}
        </p>
        <Button @click="isDeleteMode = !isDeleteMode" variant="outline">
          {{ !isDeleteMode ? "Return Books" : "Cancel" }}
        </Button>
      </div>
      <div v-if="studentRentalList.length > 0">
        <ul class="overflow-auto h-[30rem] flex flex-col gap-2">
          <Card
            v-for="rental in studentRentalList"
            :key="rental.id"
            class="relative p-3 flex flex-col gap-2"
          >
            <div
              @click="handleRentalReturn(rental.id)"
              v-if="isDeleteMode"
              class="absolute font-bold top-0 left-0 w-full h-full bg-accent opacity-80 flex items-center justify-center"
            >
              Return Book
            </div>
            <p
              v-if="rental.expand?.book_instance?.expand.book"
              class="font-bold"
            >
              {{ rental.expand?.book_instance?.expand.book.title }}
            </p>
            <div class="flex gap-2 text-xs items-baseline">
              <p
                v-if="rental.expand?.book_instance?.expand.book"
                class="bg-primary rounded-full px-2 py-0.5 text-primary-foreground"
              >
                {{ rental.expand?.book_instance?.book_code }}
              </p>
              <p v-if="rental.expand?.book_instance?.expand.book">
                {{ rental.expand?.book_instance?.expand.book.isbn }}
              </p>
            </div>

            <div class="flex justify-between items-baseline">
              <p class="text-xs">
                Borrowed {{ dateToLocaleString(rental.created) }}
              </p>
            </div>
          </Card>
        </ul>
      </div>
      <div v-else>This student has not rented any books.</div>
    </div>
  </div>
</template>

<style></style>
