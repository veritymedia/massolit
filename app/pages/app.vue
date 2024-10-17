<script lang="ts" setup>
// const pb = usePocketbase();
definePageMeta({
  // middleware: ["not-authed-guard"],
});

const pb = usePocketbase();
// const mb = useManagebac();

const allStudents = ref();
async function searchManagebacStudents() {
  try {
    const res = await pb.send("/managebac/students?q=", {});
    allStudents.value = res.students;
    console.log(res);
  } catch (err) {
    console.log(err);
  }
}

const [isOpen, closeDialog] = useDialogState();

const qrResult = ref();
function handleQrResult(result: string) {
  closeDialog();
  qrResult.value = result;
  // alert("QR:" + result);
}

async function checkIfRented(scannedCode: string) {
  try {
    const res = await pb
      .collection("rentals")
      .getFirstListItem(`book_instance.book_code="${scannedCode}"`, {
        expand: "book_instance",
      });
    if (res.id) {
      // code good and is registered.
      // fetch ManageBac user.
    }
  } catch (err) {
    console.log(err);
  }
}

onMounted(() => {
  checkIfRented("GPSYCH-1");
});
</script>

<template>
  <div class="h-screen w-screen flex-col p-6 flex gap-10">
    QR: {{ qrResult }}
    <Dialog v-model:open="isOpen" class="max-h-min">
      <DialogTrigger as-child>
        <Button class="w-full">Scan Book Code</Button>
      </DialogTrigger>
      <DialogContent class="h-full flex flex-col justify-between">
        <DialogHeader>
          <DialogTitle>Register Book</DialogTitle>
          <DialogDescription> Scan book QR code. </DialogDescription>
        </DialogHeader>
        <ScannerQrScanner @qrResult="handleQrResult" />

        <DialogFooter class="w-full flex items-center">
          <Button @click="closeDialog"> Next</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
