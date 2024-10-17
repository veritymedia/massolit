<script lang="ts" setup>
// const pb = usePocketbase();
definePageMeta({
  // middleware: ["not-authed-guard"],
});
</script>

<template>
  <div class="h-screen w-screen flex-col p-6 flex gap-10">
    <!-- <Card class="w-full h-32 xl:w-32 xl:h-full rounded">nav</Card> -->

    <Dialog class="max-h-min">
      <DialogTrigger as-child>
        <Button class="w-full">Scan Book Code</Button>
      </DialogTrigger>
      <DialogContent class="h-full flex flex-col justify-between">
        <DialogHeader>
          <DialogTitle>Register Book</DialogTitle>
          <DialogDescription> Scan book QR code. </DialogDescription>
        </DialogHeader>
        <ScannerQrScanner />

        <DialogFooter class="w-full flex items-center">
          <!-- <Button type="submit"> Next</Button> -->
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
