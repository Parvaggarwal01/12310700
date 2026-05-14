import axios from 'axios';

const registrationData = {
  email: "parvaggarwal130@gmail.com",
  name: "Parv Aggarwal",
  mobileNo: "9050547124",
  githubUsername: "parvaggarwal01",
  rollNo: "12310700",
  accessCode: "TRvZWq",
};

async function register() {
  try {
    const res = await axios.post(
      "http://4.224.186.213/evaluation-service/register",
      registrationData,
    );
    console.log("Registration Successful! SAVE THESE CREDENTIALS:");
    console.log(res.data);
  } catch (err) {
    console.error("Failed:", err.response?.data || err.message);
  }
}

register();
