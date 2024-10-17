<template>
  <div class="flex items-center justify-center">
    <div
      id="qr-scanner-container"
      :style="scannerContainerSize"
      class="w-full aspect-square"
    >
      <div id="qr-scanner-box" class="w-full h-full border" ref="qrScannerBox">
        <div id="scanner-overlay"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from "vue";
import { Html5Qrcode } from "html5-qrcode";
import { ClientResponseError } from "pocketbase";

const { width: windowWidth, height } = useWindowSize();
const scannerContainerSize = computed(() => {
  return { width: windowWidth.value - 100 };
});

const qrScannerBox = ref(null);
let html5QrCode: Html5Qrcode | null = null;

const scanResultText = ref("");

const createScanner = () => {
  html5QrCode = new Html5Qrcode("qr-scanner-box");
};

function startScanner() {
  if (!html5QrCode) {
    return;
  }
  html5QrCode
    .start(
      { facingMode: "environment" }, // Use back camera
      {
        fps: 20, // Frames per second
        aspectRatio: 1,

        // qrbox: { width: 200, height: 200 }, // Define scanner box size
      },
      (decodedText, decodedResult) => {
        console.log(`QR Code detected: ${decodedText}`);
        handleScannedCode(decodedText);
        scanResultText.value = decodedText;
      },
      (errorMessage) => {
        console.log(`QR Code scan error: ${errorMessage}`);
      },
    )
    .catch((err) => {
      console.error(`Unable to start the QR Code scanner: ${err}`);
    });
}

const stopScanner = () => {
  if (html5QrCode) {
    html5QrCode
      .stop()
      .then(() => {
        console.log("QR Code scanning stopped.");
      })
      .catch((err) => {
        console.error(`Unable to stop scanning: ${err}`);
      });
  }
};

// Start scanner when component is mounted
onMounted(() => {
  nextTick(() => {
    createScanner();
    startScanner();
  });
});

// Stop scanner when component is destroyed
onBeforeUnmount(() => {
  stopScanner();
});

const pb = usePocketbase();

async function handleScannedCode(qrResultText: string) {
  stopScanner();
  try {
    const bookInstanceResult = await pb
      .collection("book_inventory")
      .getFirstListItem("");
    // check if book exists in the database by ISBN.
    // If book does not exist in DB, fetch ISBN data from external api and ask to add it.
    // if adding book, navigate to book add screen.
    // if book in DB, then ask to add this code.
    // run codeAdd function/ call to DB.
    // once code added, move to studentBook register page
  } catch (error) {
    if (error instanceof ClientResponseError && error.status === 404) {
    }
  }

  // check if code exists in database
  // if yes, check if it is withdrawn by someone
  //    if yes return user, navigateTo "rental" page which shows book and user who took it
  //    if no, navigateTo page to search for managebac user.
  // if no:
  //    ask if should create book in database.
  //        if yes, add book to db, go back to scan screen.
  //        if no, go back to scan screen.
}
</script>

<style scoped>
#scanner-overlay {
  position: absolute;
  top: 0;
  left: 0;
  border: 2px;
  width: 100%;
  height: 100%;
  background: rgba(5, 5, 5, 0.7); /* Overlay effect */
  pointer-events: none;
}
</style>
