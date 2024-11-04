<script lang="ts" setup>

const pb = usePocketbase()
const isLoading = ref(false)
const managebacResult = ref();
const selectedStudent = ref()
const studentSearchTerm = ref()


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
        // studentListOpen.value = true;
        console.log(studentSearchTerm.value);
        await searchManagebacStudents();
    }, 750);
};


function handleStudentSelect(student: any) {
    console.log("Selected student: ", student.email);
    selectedStudent.value = student;
    studentSearchTerm.value = "";
    managebacResult.value = {};
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


async function getRentalsByStudentId() { }

</script>

<template>
    <div>
        <h2>Find rented books by student</h2>
        <div class="w-full flex items-center items-evenly border border-muted rounded-lg">
            <Input @input="handleInput" class="border-0" placeholder="Search for student" v-model="studentSearchTerm" />
            <Icon v-if="isLoading" name="line-md:loading-loop" class="mx-2" />
        </div>
        <Card v-if="studentList.length !== 0 && studentSearchTerm.length !== 0"
            class="mt-2 absolute top-10 bg-background z-50 w-full overflow-scroll">
            <ul class="flex flex-col">
                <li class="hover:bg-muted cursor-pointer first:mt-3 last:mb-3 py-0.5 px-3"
                    @click="handleStudentSelect(student)" v-for="student in studentList" :key="student.id">
                    {{ student.first_name }} {{ student.last_name }}
                </li>
            </ul>
        </Card>

    </div>
</template>