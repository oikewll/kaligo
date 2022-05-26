/**
 * 该文件用于nuxt ssr项目，nuxt build 会产生用于nodejs的文件，在.nuxt/dist目录下，需要脚本手动操作
 * 部署目录，规范是 dist/webfe_{xxx}
 * 1. 清空目录
 * 2. 建立路径
 * 3. 复制deploy所需文件
 * 4. 进入到部署目录执行npm install --production
 * 
 * @author Dawsonliu
 * @date 2021/11/09
 */
const pkg = require("../package.json");
const gulp = require("gulp");
const del = require("del");
const spawn = require("child_process").spawn;
const path = require("path");
const { createPath } = require("./tools");

const ROOT_PATH = path.resolve(__dirname, '../');
const APP_PATH = `${ROOT_PATH}/dist/${pkg.name}`;
const DEST_PATH = path.resolve(__dirname, APP_PATH);

const cleanAndBuild = (done) => {
    try {
        del.sync([`${DEST_PATH}/**`, `!${DEST_PATH}`], {
            force: true,
            dot: true,
        });
        console.log("🎉Clean dir");

        createPath(APP_PATH);
        console.log("🎉Create dest path");
        done();
    } catch (err) {
        throw new Error(`🙅‍♂️cleanAndBuild error： ${err}`);
    }
}

const copyDeploy = async (done) => {
    try {
        await Promise.all([
            copyPromise(`${ROOT_PATH}/.nuxt/**`, `${DEST_PATH}/.nuxt`),
            copyPromise(`${ROOT_PATH}/store/**`, `${DEST_PATH}/store`),
            copyPromise(`${ROOT_PATH}/static/**`, `${DEST_PATH}/static`),
            copyPromise(`${ROOT_PATH}/config/**`, `${DEST_PATH}/config`),
            copyPromise(`${ROOT_PATH}/gulpfile.js/**`, `${DEST_PATH}/gulpfile.js`),
            copyPromise(`${ROOT_PATH}/deploy/index.js`, `${DEST_PATH}`),
            copyPromise(`${ROOT_PATH}/nuxt.config.js`, `${DEST_PATH}`),
            copyPromise(`${ROOT_PATH}/local.config.js`, `${DEST_PATH}`),
            copyPromise(`${ROOT_PATH}/.iconfont`, `${DEST_PATH}`),
            copyPromise(`${ROOT_PATH}/.npmrc`, `${DEST_PATH}`),
            copyPromise(`${ROOT_PATH}/package.json`, `${DEST_PATH}`),
        ]);
        console.log("🎉Copy done");
        done();
    } catch (err) {
        throw new Error(`🙅‍♂️copyDeploy error： ${err}`);
    }
}

const installDeploy = (done) => {
    spawn("npm", ["install", "--production"], {
        cwd: DEST_PATH,
        stdio: "inherit",
    })
        .on("error", (err) => {
            throw new Error(`🙅‍♂️installDeploy error： ${err}`);
        })
        .on("close", () => {
            console.log("🙅installDeploy done!!")
            done();
        });
}

// 复制文件promise
const copyPromise = (src, dest) => {
    return new Promise((resolve, reject) => {
        const stream = gulp.src(src).pipe(gulp.dest(dest));
        stream.on("end", function () {
            console.log("😔copy", src);
            resolve();
        });
        stream.on("error", reject);
    });
}

exports.cleanAndBuild = cleanAndBuild;
exports.copyDeploy = copyDeploy;
exports.installDeploy = installDeploy;

gulp.task('build', gulp.series(
    cleanAndBuild,
    copyDeploy,
    installDeploy
));